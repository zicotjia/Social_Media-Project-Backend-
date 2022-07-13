package middleware

import (
	helper "backend/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

		c.Next()

	}
}

//var ctx context.Context
//
//var col mongo.Collection
//
//func GetAllTask(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	payload := getAllTask()
//	json.NewEncoder(w).Encode(payload)
//}
//
//func getAllTask() []bson.M {
//	cursor, err := col.Find(context.TODO(), bson.M{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	var results []bson.M
//	for cursor.Next(ctx) {
//		var result bson.M
//		if err = cursor.Decode(&result); err != nil {
//			log.Fatal(err)
//		}
//		results = append(results, result)
//	}
//	fmt.Println(results)
//	return results
//
//}
//
//func main() {
//
//	var uri = "mongodb+srv://123456:123456zico@cluster0.sbbtypy.mongodb.net/?retryWrites=true&w=majority"
//	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
//	if err != nil {
//		log.Fatal(err)
//	}
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	err = client.Connect(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Disconnect(ctx)
//
//	//quickstartDatabase := client.Database("quickstart")
//	//podcastsCollection := quickstartDatabase.Collection("podcasts")
//	//episodesCollection := quickstartDatabase.Collection("episodes")
//	//podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
//	//	{Key: "title", Value: "The Polyglot Developer Podcast"},
//	//	{Key: "author", Value: "Nic Raboy"},
//	//})
//
//	//Creating new entries
//	//episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
//	//	bson.D{
//	//		{"podcast", podcastResult.InsertedID},
//	//		{"title", "GraphQL for API Development"},
//	//		{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
//	//		{"duration", 25},
//	//	},
//	//	bson.D{
//	//		{"podcast", podcastResult.InsertedID},
//	//		{"title", "Progressive Web Application Development"},
//	//		{"description", "Learn about PWA development with Tara Manicsic."},
//	//		{"duration", 32},
//	//	},
//	//})
//	//fmt.Println("inserted" + fmt.Sprintf("%v", podcastResult.InsertedID))
//	//fmt.Println("inserted" + fmt.Sprintf("%v", episodeResult.InsertedIDs))
//
//	//reading only one entry (find the first entry)
//	//var podcast bson.M
//	//if err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
//	//	log.Fatal(err)
//	//}
//	//fmt.Println(podcast)
//
//	//read multiple entries that passed a filter
//	//var episode []bson.M
//	//filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//if err = filterCursor.All(ctx, &episode); err != nil {
//	//	log.Fatal(err)
//	//}
//
//	//query and then sort
//	//opts := options.Find()
//	//opts.SetSort(bson.D{{"duration", -1}})
//	//sortCursor, err := episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//var episodesSorted []bson.M
//	//if err = sortCursor.All(ctx, &episodesSorted); err != nil {
//	//	log.Fatal(err)
//	//}
//	//fmt.Println(episodesSorted)
//
//	podcastsCollection := client.Database("quickstart").Collection("podcasts")
//	// if of entry to be updated
//	id, _ := primitive.ObjectIDFromHex("62a2e49b4f85fa3405a49709")
//	result, err := podcastsCollection.UpdateOne(
//		ctx,
//		bson.M{"_id": id},
//		bson.D{
//			{"$set", bson.D{{"author", "Nic Raboy"}}},
//			{"$currentDate", bson.D{{"dateModified", true}}},
//		},
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
//}
