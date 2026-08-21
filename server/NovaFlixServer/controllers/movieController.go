package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/database"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/models"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/googleai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenConnection("movies")
var rankingCollection *mongo.Collection = database.OpenConnection("rankings")

var validate=validator.New()
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch movies"})
				return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode movies"})
			return

		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}
		var movie models.Movie

		err:= movieCollection.FindOne(ctx,bson.M{"imdb_id":movieID}).Decode(&movie)
		if err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"Movie not Found"})
			return


		}

		c.JSON(http.StatusOK,movie)
	}
	
}

func AddMovie() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var movie models.Movie

		if err:=c.ShouldBindJSON(&movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Movie Details are required to add Movie"})
			return
		}

		if err:=validate.Struct(movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Validation Failed","details":err.Error()})
			return 
		}
		result,err:=movieCollection.InsertOne(ctx,movie)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to add Movie"})
			return
		}

		c.JSON(http.StatusCreated,result)

	}
}


func AdminReviewUpdate() gin.HandlerFunc{
	return  func(c *gin.Context) {
		role,err:=utils.GetUserRoleFromContext(c)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Role not found in context"})
			return
		}
		if role!="ADMIN"{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"You are not authorized for this functionality"})
			return
		}
		movieId:=c.Param("imdb_id")
		if movieId==""{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Movie Id Required"})
			return 
		}
		var req struct{
			AdminReview string `json:"admin_review"`
		}
		var res struct{
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err:=c.ShouldBind(&req);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request body"})
			return 
		}

		sentiment,rankVal,err:=GetReviewRanking(req.AdminReview)

		if err!=nil{
			log.Println("GetReviewRanking error:", err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error getting review ranking"})
			return
		}
		filter:=bson.M{"imdb_id":movieId}
		update:=bson.M{
			"$set":bson.M{
				"admin_review":req.AdminReview,
				"ranking":bson.M{
					"ranking_value":rankVal,
					"ranking_name":sentiment,
				},
			},
		}

		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		result,err:=movieCollection.UpdateOne(ctx,filter,update)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Error updating movie"})
			return 
		}

		if result.MatchedCount==0{
			c.JSON(http.StatusNotFound,gin.H{"error":"Movie not found"})
			return
		}
		res.AdminReview=req.AdminReview
		res.RankingName=sentiment
		c.JSON(http.StatusOK,res)
	}
	
}


func GetReviewRanking(admin_review string)(string,int,error){
	rankings,err:=GetRankings()

	if err!=nil{
		return "" ,0,nil
	}
	sentimentDelimited:=""
	for _,ranking:=range rankings{
		if ranking.RankingValue!=999{
			sentimentDelimited=sentimentDelimited+ranking.RankingName+","
		}

	}

	sentimentDelimited=strings.Trim(sentimentDelimited,",")

	err=godotenv.Load(".env")

	if err!=nil{
		log.Println("Warning: .env file not found")
	}

	GeminiKey:=os.Getenv("GEMINI_API_KEY")

	if GeminiKey==""{
		return "",0,errors.New("Could not read GEMINI_API_KEY")
	}

	llm,err:= googleai.New(context.Background(),googleai.WithAPIKey(GeminiKey), googleai.WithDefaultModel("gemini-3.6-flash"))

	if err!=nil{
		return "",0,err
	}

	prompt_Template:=os.Getenv("BASE_PROMPT_TEMPLATE")
	if prompt_Template==""{
		return "",0,errors.New("Could not read BASE_PROMPT")

	}
	base_Prompt:=strings.Replace(prompt_Template,"{rankings}",sentimentDelimited,1)

	response,err:=llm.Call(context.Background(),base_Prompt+admin_review)

	if err!=nil{
		return "",0,err
	}
	rankVal:=0
	for _,ranking:=range rankings{
		if ranking.RankingName==response{
			rankVal=ranking.RankingValue
			break
		}

	} 
	return response,rankVal,nil
	

}

func GetRankings()([]models.Ranking,error){
	var rankings []models.Ranking

	var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)

	defer cancel()

	cursor,err:= rankingCollection.Find(ctx,bson.M{})

	if err!=nil{
		return nil,err
	}
	defer cursor.Close(ctx)

	if err:=cursor.All(ctx,&rankings);err!=nil{
		return nil,err
	}
	return rankings,nil
}



func GetRecommendedMovies() gin.HandlerFunc{
	return func(c *gin.Context) {
		userId,err:=utils.GetUserIdFromContext(c)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"User Id not found in context"})
			return 
		}
		favourite_genres,err:=GetUsersFavouriteGenres(userId)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return 
		}
		err=godotenv.Load(".env")
			if err!=nil{
				log.Println("Warning: .enc file not found")
			}
			
		var recommendedMovieLimitVal int64=5
		recommendedMovieLimitStr:=os.Getenv("RECOMMENDED_MOVIE_VALUE") 
		if recommendedMovieLimitStr!=""{
			recommendedMovieLimitVal,_=strconv.ParseInt(recommendedMovieLimitStr,10,64)
		}

		findOptions:= options.Find()

		findOptions.SetSort(bson.D{{Key:"ranking.ranking_value",Value:1}})
		findOptions.SetLimit(recommendedMovieLimitVal)
		filter:=bson.M{"genre.genre_name":bson.M{"$in":favourite_genres}}

		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
	defer cancel()
	cursor,err:=movieCollection.Find(ctx,filter,findOptions)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Error fetching recommended movies"})
		return
	}
	defer cursor.Close(ctx)

	var recommendedMovies []models.Movie
	if err:=cursor.All(ctx,&recommendedMovies);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK,recommendedMovies)


	}
}





func GetUsersFavouriteGenres(userId string)([]string,error){
	var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
	defer cancel()

	filter:=bson.M{"user_id":userId}

	projection:=bson.M{
		"favourite_genres":1,
		"_id":0,
	}
	opts:=options.FindOne().SetProjection(projection)

	var result bson.M
	err:=userCollection.FindOne(ctx,filter,opts).Decode(&result)
	if err!=nil{
		if err==mongo.ErrNoDocuments{
			return []string{},nil
		}

	}
	favGenresArray,ok:=result["favourite_genres"].(bson.A)
	if !ok{
		return []string{},errors.New("Unable to retrieve favourite geenres of user")
	}
	var genreNames []string

	for _,item:=range favGenresArray{
		if genreMap,ok:=item.(bson.D);ok{
			for _,elem:=range genreMap{
				if elem.Key=="genre_name"{
					if name,ok:=elem.Value.(string);ok{
						genreNames=append(genreNames, name)
					}
				}
			}
		}
	}
	return genreNames,nil
}
