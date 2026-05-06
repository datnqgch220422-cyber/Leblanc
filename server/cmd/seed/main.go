package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()
	db.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := ensureUserIndexes(ctx); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}
	if err := seedDrinks(ctx); err != nil {
		log.Fatalf("seed drinks: %v", err)
	}

	log.Println("Database ready.")
}

func ensureUserIndexes(ctx context.Context) error {
	coll := db.DB.Collection("users")
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "nameLower", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "emailLower", Value: 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := coll.Indexes().CreateMany(ctx, models)
	return err
}

func seedDrinks(ctx context.Context) error {
	coll := db.DB.Collection("drinks")
	if _, err := coll.DeleteMany(ctx, bson.D{}); err != nil {
		return err
	}

	assetBase := strings.TrimRight(os.Getenv("DRINK_IMAGE_BASE"), "/")
	if assetBase == "" {
		assetBase = "https://le-blanc-web.vercel.app/drinks"
	}

	drinks := []models.Drink{
		// Day menu (12)
		{
			ID:         primitive.NewObjectID(),
			Name:       "Vietnamese iced milk coffee",
			Price:      29000,
			Tags:       []string{"day", "coffee", "phin", "iced", "milk"},
			Caffeine:   "high",
			Temp:       "iced",
			Sweetness:  4,
			ColorTone:  "warm",

			Image:      assetBase + "/cafe-sua-da.png",
			Desc:       "Robusta phin from Buon Ma Thuot, condensed milk, ice; bold and creamy Saigon style.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Orange lemongrass americano",
			Price:      39000,
			Tags:       []string{"day", "coffee", "espresso", "citrus", "iced"},
			Caffeine:   "med",
			Temp:       "iced",
			Sweetness:  2,
			ColorTone:  "cool",

			Image:      assetBase + "/americano-cam-sa.png",
			Desc:       "Da Lat Arabica espresso with fresh orange juice and smashed lemongrass; bright citrus aroma.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Oat milk latte",
			Price:      45000,
			Tags:       []string{"day", "coffee", "espresso", "latte", "oat"},
			Caffeine:   "med",
			Temp:       "hot",
			Sweetness:  3,
			ColorTone:  "neutral",

			Image:      assetBase + "/latte-sua-yen-mach.png",
			Desc:       "Arabica espresso with unsweetened oat milk, fine foam; dairy-free friendly.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Coconut iced coffee",
			Price:      49000,
			Tags:       []string{"day", "coffee", "coconut", "blended", "iced"},
			Caffeine:   "med",
			Temp:       "iced",
			Sweetness:  4,
			ColorTone:  "warm",

			Image:      assetBase + "/cafe-dua.png",
			Desc:       "Arabica-Robusta espresso blend with light coconut cream and low-sugar milk, blended with ice.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Peach orange lemongrass tea",
			Price:      39000,
			Tags:       []string{"day", "tea", "fruit", "iced"},
			Caffeine:   "low",
			Temp:       "iced",
			Sweetness:  3,
			ColorTone:  "warm",

			Image:      assetBase + "/tra-dao-cam-sa.png",
			Desc:       "Cold-brew Ceylon black tea with peach syrup, orange slices, lemongrass, and peach chunks.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Ginger honey tea",
			Price:      32000,
			Tags:       []string{"day", "tea", "ginger", "hot"},
			Caffeine:   "low",
			Temp:       "hot",
			Sweetness:  2,
			ColorTone:  "warm",

			Image:      assetBase + "/tra-gung-mat-ong.png",
			Desc:       "Light green tea, fresh ginger slices, and forest honey; soothing for cool mornings.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Orange carrot juice",
			Price:      39000,
			Tags:       []string{"day", "juice", "fresh", "cold"},
			Caffeine:   "none",
			Temp:       "iced",
			Sweetness:  3,
			ColorTone:  "neutral",

			Image:      assetBase + "/nuoc-ep-cam-ca-rot.png",
			Desc:       "Fresh orange and carrot juice, no syrup; sweetness adjustable on request.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Banana oat smoothie",
			Price:      45000,
			Tags:       []string{"day", "smoothie", "healthy"},
			Caffeine:   "none",
			Temp:       "iced",
			Sweetness:  3,
			ColorTone:  "neutral",

			Image:      assetBase + "/sinh-to-chuoi-yen-mach.png",
			Desc:       "Ripe banana, rolled oats, almond milk, touch of honey; filling and gym-friendly.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Passionfruit yogurt frappe",
			Price:      39000,
			Tags:       []string{"day", "yogurt", "blended"},
			Caffeine:   "none",
			Temp:       "iced",
			Sweetness:  4,
			ColorTone:  "cool",

			Image:      assetBase + "/sua-chua-chanh-day.png",
			Desc:       "Fermented yogurt with fresh passionfruit sauce, blended with ice; tangy-sweet and creamy.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Matcha latte",
			Price:      49000,
			Tags:       []string{"day", "tea", "matcha", "latte", "hot"},
			Caffeine:   "low",
			Temp:       "hot",
			Sweetness:  3,
			ColorTone:  "neutral",

			Image:      assetBase + "/matcha-latte.png",
			Desc:       "Japanese matcha whisked with fresh milk, steamed; aromatic with gentle bitterness.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Hot cacao with cinnamon",
			Price:      45000,
			Tags:       []string{"day", "cacao", "hot"},
			Caffeine:   "low",
			Temp:       "hot",
			Sweetness:  4,
			ColorTone:  "warm",

			Image:      assetBase + "/ca-cao-que-nong.png",
			Desc:       "Natural cacao powder, fresh milk, optional condensed milk, honey or syrup, plus cinnamon dust.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Warm apple cinnamon honey tea",
			Price:      39000,
			Tags:       []string{"day", "tea", "apple", "cinnamon", "hot"},
			Caffeine:   "low",
			Temp:       "hot",
			Sweetness:  3,
			ColorTone:  "warm",

			Image:      assetBase + "/tra-tao-que-mat-ong-nong.png",
			Desc:       "Light tea with fresh apple slices, cinnamon stick, honey; optional squeeze of lime.",
		},
		// Night menu (10)
		{
			ID:         primitive.NewObjectID(),
			Name:       "Midnight Coffee",
			Price:      129000,
			Tags:       []string{"night", "cocktail", "coffee", "signature"},
			Caffeine:   "low",
			Temp:       "cold",
			Sweetness:  4,
			ColorTone:  "warm",

			Image:      assetBase + "/midnight-coffee.png",
			Desc:       "Coffee cocktail with rum or whisky and espresso; warm, boozy coffee sweetness.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Citrus Gin & Tonic",
			Price:      99000,
			Tags:       []string{"night", "cocktail", "gin", "citrus"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  2,
			ColorTone:  "cool",

			Image:      assetBase + "/gin-tonic-chanh-buoi.png",
			Desc:       "Gin and tonic with lemon and grapefruit peel; crisp, aromatic, lightly fizzy.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Classic Mint Mojito",
			Price:      109000,
			Tags:       []string{"night", "cocktail", "rum", "mint"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  3,
			ColorTone:  "neutral",

			Image:      assetBase + "/mojito-bac-ha-co-dien.png",
			Desc:       "White rum, fresh mint, lime, soda; cool and gently sweet.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Rum & Coke with lime",
			Price:      89000,
			Tags:       []string{"night", "cocktail", "rum", "cola"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  3,
			ColorTone:  "warm",

			Image:      assetBase + "/rum-coke-chanh.png",
			Desc:       "White rum and cola with lime; familiar, bubbly, easy-drinking.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Whisky Highball",
			Price:      129000,
			Tags:       []string{"night", "cocktail", "whisky", "highball"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  2,
			ColorTone:  "neutral",

			Image:      assetBase + "/whisky-highball.png",
			Desc:       "Whisky with soda, subtle fizz; mellow malt aroma, not harsh.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "House Red Wine (glass)",
			Price:      119000,
			Tags:       []string{"night", "wine", "red"},
			Caffeine:   "none",
			Temp:       "room",
			Sweetness:  2,
			ColorTone:  "warm",

			Image:      assetBase + "/house-red.png",
			Desc:       "Fruity red wine, slightly tart; easy, dinner-friendly.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "House White Wine (glass)",
			Price:      119000,
			Tags:       []string{"night", "wine", "white"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  2,
			ColorTone:  "cool",

			Image:      assetBase + "/house-white.png",
			Desc:       "Bright white wine with mild apple, pear, and citrus notes.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Classic lager beer",
			Price:      49000,
			Tags:       []string{"night", "beer", "lager"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  1,
			ColorTone:  "neutral",

			Image:      assetBase + "/beer.png",
			Desc:       "Cold lager, lightly bitter, crisp and easy.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Craft beer (rotating taps)",
			Price:      89000,
			Tags:       []string{"night", "beer", "craft"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  1,
			ColorTone:  "warm",

			Image:      assetBase + "/crafted-beer.png",
			Desc:       "Pale ale, IPA, or wheat rotates; mildly bitter, hoppy or fruity depending on tap.",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "Plum liqueur (umeshu)",
			Price:      89000,
			Tags:       []string{"night", "liqueur", "plum", "sweet"},
			Caffeine:   "none",
			Temp:       "cold",
			Sweetness:  4,
			ColorTone:  "warm",

			Image:      assetBase + "/plum-liqueur.png",
			Desc:       "Japanese-style plum liqueur; sweet fruit aroma, guest-friendly.",
		},
	}

	docs := make([]any, len(drinks))
	for i, d := range drinks {
		docs[i] = d
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return err
	}

	log.Printf("seeded %d drinks (day + night)", len(drinks))
	return nil
}
