package repo

import (
	"context"
	"database-example/model"
	"fmt"
	"log"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*FollowerRepository, error) {
	uri := "bolt://localhost:7687"
	user := "neo4j"
	pass := "neo4jteam25"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	logger.Println("Connected to Neo4j database")
	return &FollowerRepository{
		driver: driver,
		logger: logger,
	}, nil
}

func (fr *FollowerRepository) CheckConnection() {
	ctx := context.Background()
	err := fr.driver.VerifyConnectivity(ctx)
	if err != nil {
		fr.logger.Panic(err)
		return
	}

	fr.logger.Printf(`Neo4J server address: %s`, fr.driver.Target().Host)
}

func (mr *FollowerRepository) CreateUser(user *model.User) error {
	ctx := context.Background()
	session := mr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	savedPerson, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (u:User {id: $id, username: $username}) RETURN u.username + ', from node ' + id(u)",
				map[string]interface{}{"id": user.ID, "username": user.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		mr.logger.Println("Error inserting Person:", err)
		return err
	}
	mr.logger.Println(savedPerson.(string))
	return nil
}

func (fr *FollowerRepository) CreateFollowers(user *model.User, followingUser *model.User) error {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// Ispis informacija o korisnicima
	fmt.Println("User ID:", user.ID, "Username:", user.Username)
	fmt.Println("Following User ID:", followingUser.ID, "Username:", followingUser.Username)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				`
				MATCH (u:User {id: $userID}), (f:User {id: $followingUserID})
				CREATE (u)-[:FOLLOWS]->(f)
				RETURN u, f
				`,
				map[string]interface{}{
					"userID":          user.ID,
					"followingUserID": followingUser.ID,
				})
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

	if err != nil {
		fr.logger.Println("Error creating followers relationship:", err)
		return err
	}

	// Kreiranje Followers objekta
	followers := &model.Followers{
		UserId:            strconv.Itoa(user.ID),
		Username:          user.Username,
		FollowingUserId:   strconv.Itoa(followingUser.ID),
		FollowingUsername: followingUser.Username,
	}

	// Ispis Followers objekta
	fmt.Println("Followers object:", followers)

	// Upisivanje Followers objekta u Neo4j bazu
	_, err = session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			_, err := transaction.Run(ctx,
				`
				CREATE (follower:Followers {
					userId: $userId,
					username: $username,
					followingUserId: $followingUserId,
					followingUsername: $followingUsername
				})
				RETURN follower
				`,
				map[string]interface{}{
					"userId":            followers.UserId,
					"username":          followers.Username,
					"followingUserId":   followers.FollowingUserId,
					"followingUsername": followers.FollowingUsername,
				})
			if err != nil {
				return nil, err
			}

			return nil, nil
		})

	if err != nil {
		fr.logger.Println("Error creating Followers node:", err)
		return err
	}

	return nil
}

func (fr *FollowerRepository) GetAllFollowers() ([]*model.Followers, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.Run(
		ctx,
		`MATCH (follower:Followers)
        RETURN follower.userId as userId, follower.username as username, follower.followingUserId as followingUserId, follower.followingUsername as followingUsername`,
		nil,
	)
	if err != nil {
		fr.logger.Println("Error getting followers:", err)
		return nil, err
	}

	var followers []*model.Followers
	for result.Next(ctx) {
		record := result.Record()
		userId, _ := record.Get("userId")
		username, _ := record.Get("username")
		followingUserId, _ := record.Get("followingUserId")
		followingUsername, _ := record.Get("followingUsername")
		follower := &model.Followers{
			UserId:            userId.(string),
			Username:          username.(string),
			FollowingUserId:   followingUserId.(string),
			FollowingUsername: followingUsername.(string),
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

func (fr *FollowerRepository) GetRecommendations(userID string) ([]*model.User, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.Run(
		ctx,
		`
        MATCH (u:User)-[:FOLLOWS]->(f:User)<-[:FOLLOWS]-(ff:User)
        WHERE u.id = $userID
        RETURN ff.id as id, ff.username as username
        `,
		map[string]interface{}{"userID": userID})
	if err != nil {
		fr.logger.Println("Error getting recommendations:", err)
		return nil, err
	}

	var recommendations []*model.User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		username, _ := record.Get("username")
		userIDInt, _ := strconv.Atoi(id.(string))
		user := &model.User{
			ID:       userIDInt,
			Username: username.(string),
		}
		recommendations = append(recommendations, user)
	}

	return recommendations, nil
}

func (fr *FollowerRepository) GetFollowings(userID string) ([]*model.User, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// Konvertuj userID u int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fr.logger.Println("Error converting userID to int:", err)
		return nil, err
	}

	// Formatiranje upita sa stvarnom vrednošću userID
	query := fmt.Sprintf(`
    MATCH (u:User)-[:FOLLOWS]->(f:User)
    WHERE f.id = %d
    RETURN u
    `, userIDInt)

	// Dodajte log koji će ispisati formatiran upit
	fmt.Println("Query:", query)

	fmt.Println("UserID:", userID)
	result, err := session.Run(
		ctx,
		query,
		nil, // Ne koristimo map[string]interface{} jer smo direktno uneli vrednost u upit
	)
	if err != nil {
		fr.logger.Println("Error getting followings:", err)
		return nil, err
	}

	var followers []*model.User
	// Ispisivanje rezultata pre iteracije
	fmt.Println("Raw results from the database:")
	for result.Next(ctx) {
		record := result.Record()
		userNode, found := record.Get("u")
		if !found {
			continue
		}
		fmt.Print(userNode) //PRAVO PROBLEM OKO PARSIRANJA.................NE ZNAM ISKRENO VISE

		// Parsiranje rezultata
		props, ok := userNode.(map[string]interface{})
		if !ok {
			fmt.Println("Error parsing user node properties")
			continue
		}

		userIDValue, ok := props["id"].(int64)
		if !ok {
			fmt.Println("Error getting user ID")
			continue
		}

		username, ok := props["username"].(string)
		if !ok {
			fmt.Println("Error getting username")
			continue
		}

		// Kreiranje objekta User
		userObj := &model.User{
			ID:       int(userIDValue),
			Username: username,
		}
		followers = append(followers, userObj)

		// Iteriranje kroz rezultate
		// for result.Next(ctx) {
		// 	record := result.Record()

		// 	// Izvlačenje vrednosti iz rezultata
		// 	id, ok := record.Get("u.id")
		// 	if !ok {
		// 		fr.logger.Println("Error getting user ID")
		// 		continue
		// 	}

		// 	username, ok := record.Get("u.username")
		// 	if !ok {
		// 		fr.logger.Println("Error getting username")
		// 		continue
		// 	}

		// 	// Konvertovanje ID-ja u int
		// 	userIDInt, err := strconv.Atoi(id.(string))
		// 	if err != nil {
		// 		fr.logger.Println("Error converting user ID to int:", err)
		// 		continue
		// 	}

		// 	// Kreiranje objekta User
		// 	userObj := &model.User{
		// 		ID:       userIDInt,
		// 		Username: username.(string),
		// 	}
		// 	followers = append(followers, userObj)
	}

	return followers, nil
}

func (mr *FollowerRepository) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}

func (fr FollowerRepository) GetUserById(userID int) (*model.User, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	result, err := session.Run(
		ctx,
		`
        MATCH (u:User {id: $userID})
        RETURN u.username as username
        `,
		map[string]interface{}{"userID": userID})
	if err != nil {
		fr.logger.Println("Error getting user by ID:", err)
		return nil, err
	}

	if result.Next(ctx) {
		record := result.Record()
		username, _ := record.Get("username")
		user := &model.User{
			ID:       userID,
			Username: username.(string),
		}
		return user, nil
	}

	return nil, nil
}
