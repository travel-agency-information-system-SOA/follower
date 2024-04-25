package repo

import (
	"context"
	"database-example/model"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*FollowerRepository, error) {
	// uri := "bolt://localhost:7687"
	// user := "neo4j"
	// pass := "neo4jteam25"
	// auth := neo4j.BasicAuth(user, pass, "")
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
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
				MATCH (u:User {ID: $userID}), (f:User {ID: $followingUserID})
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
        MATCH (u:User {ID: $userID})
        RETURN u.Username as username
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

func (fr *FollowerRepository) GetFollowings(userID string) ([]*model.User, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	fr.logger.Println("usao u repo za get followings")

	// Konvertujemo userID u int64
	userIdInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		fr.logger.Println("Error converting userID to int64:", err)
		return nil, err
	}

	result, err := session.Run(
		ctx,
		`
        MATCH (u:User {ID: $userID})-[:FOLLOWS]->(f:User)
        RETURN f.ID as id, f.Username as username
        `,
		map[string]interface{}{"userID": userIdInt})
	if err != nil {
		fr.logger.Println("Error getting followings:", err)
		return nil, err
	}

	var followings []*model.User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		username, _ := record.Get("username")
		// Konvertujemo id u int
		idInt := int(id.(int64))
		user := &model.User{
			ID:       idInt,
			Username: username.(string),
		}
		followings = append(followings, user)
		fr.logger.Println("Found user - ID:", user.ID, "Username:", user.Username)
	}

	if len(followings) == 0 {
		fmt.Println("No followings found for user with ID:", userID)
	}
	return followings, nil
}

func (fr *FollowerRepository) GetRecommendations(userID string) ([]*model.User, error) {
	ctx := context.Background()
	session := fr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	fr.logger.Println("usao u repo za get recommmendations")

	// Konvertujemo userID u int64
	userIdInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		fr.logger.Println("Error converting userID to int64:", err)
		return nil, err
	}

	result, err := session.Run(
		ctx,
		`
		MATCH (u:User {ID: $userID})-[:FOLLOWS]->(:User)-[:FOLLOWS]->(f:User)
		WHERE f.ID <> $userID
		RETURN DISTINCT f.ID as id, f.Username as username
        `,
		map[string]interface{}{"userID": userIdInt})
	if err != nil {
		fr.logger.Println("Error getting followings:", err)
		return nil, err
	}

	var followings []*model.User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		username, _ := record.Get("username")
		// Konvertujemo id u int
		idInt := int(id.(int64))
		user := &model.User{
			ID:       idInt,
			Username: username.(string),
		}
		followings = append(followings, user)
		fr.logger.Println("Found user - ID:", user.ID, "Username:", user.Username)
	}

	if len(followings) == 0 {
		fmt.Println("No recommendations found for user with ID:", userID)
	}
	return followings, nil
}
