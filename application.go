package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoDBHosts    = "ds217898.mlab.com:17898"
	AuthDatabase    = "plydge-mvp-v01"
	AuthUserName    = "plydge-api"
	AuthPassword    = "api20!7"
	TestDatabase    = "plydge-mvp-v01"
	TitleCollection = "titles"
)

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(AuthDatabase).C(TitleCollection)

	index := mgo.Index{
		Key:        []string{"identificationnumber"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func main() {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	ensureIndex(session)

	type RegisteredOwner struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		City    string `json:"city"`
		State   string `json:"state"`
		ZipCode int    `json:"zipcode"`
		tStamp  string
	}

	type Salvage struct {
		Note string `json:"note"`
	}

	type Brand struct {
		Note string `json:"note"`
	}

	type Lien struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		City        string `json:"city"`
		State       string `json:"state"`
		ZipCode     int    `json:"zipcode"`
		Date        string `json:"date"`
		ReceiptDate string `json:"receiptdate"`
		tStamp      string
	}

	type Title struct {
		IdentificationNumber     string `json:"identificationnumber"`
		Year                     int    `json:"year"`
		Make                     string `json:"make"`
		Body                     string `json:"body"`
		WtlBhp                   int    `json:"wtlbhp"`
		VesselRegistrationNumber string `json:"vesselregistrationnumber"`
		TitleNumber              int    `json:"titlenumber"`
		PrevState                string `json:"prevstate"`
		Color                    string `json:"color"`
		PrimaryBrand             []Brand
		SecondaryBrand           []Brand
		NumberOfBrands           int    `json:"numberofbrands"`
		Use                      string `json:"use"`
		PreIssueDate             string `json:"previssuedate"`
		OdometerStatus           string `json:"odometerstatus"`
		HullMaterial             string `json:"hullmaterial"`
		Prop                     string `json:"prop"`
		DateOfIssue              string `json:"dateofissue"`
		Owner                    []RegisteredOwner
		RegisteredOwners         []RegisteredOwner
		FirstLienHolder          []Lien
		Liens                    []Lien
		Salvages                 []Salvage
		Brands                   []Brand
		tStamp                   string
	}

	lh := Lien{
		Name:        "June",
		Address:     "12306 Sawgrass CT",
		City:        "Wellington",
		State:       "FL",
		ZipCode:     33414,
		ReceiptDate: "01/12/2018",
		Date:        "01/11/2015",
	}

	t := Title{
		IdentificationNumber: "SRF32530",
		Year:                 1977,
		Make:                 "ROLL",
		Body:                 "4D",
		WtlBhp:               4995,
		VesselRegistrationNumber: "",
		TitleNumber:              62710196,
		PrevState:                "FL",
		Color:                    "MAR",
		PrimaryBrand:             []Brand{},
		SecondaryBrand:           []Brand{},
		NumberOfBrands:           0,
		Use:                      "Private",
		PreIssueDate:             "09/16/1992",
		OdometerStatus:           "EXEMPT",
		HullMaterial:             "",
		Prop:                     "",
		DateOfIssue:              "09/11/2007",
		RegisteredOwners: []RegisteredOwner{
			RegisteredOwner{
				Name:    "June",
				Address: "12306 Sawgrass CT",
				City:    "Wellington",
				State:   "FL",
				ZipCode: 33414,
			},
		},
		Owner: []RegisteredOwner{
			RegisteredOwner{
				Name:    "June",
				Address: "12306 Sawgrass CT",
				City:    "Wellington",
				State:   "FL",
				ZipCode: 33414,
			},
		},
		FirstLienHolder: []Lien{lh},
		Liens:           []Lien{lh, lh, lh, lh},
		Salvages:        []Salvage{},
		Brands:          []Brand{},
	}

	router := gin.Default()

	// GET a title detail
	router.GET("/title/:id", func(c *gin.Context) {
		var (
			title  Title
			result gin.H
		)
		id := c.Param("id")

		title = t

		s := session.Copy()

		defer s.Close()

		col := s.DB(AuthDatabase).C(TitleCollection)

		err := col.Find(bson.M{"identificationnumber": id}).One(&title)

		if err != nil {
			result = gin.H{
				"result":  "",
				"count":   id,
				"err":     err.Error(),
				"success": 0,
			}
			c.JSON(http.StatusOK, result)
		}

		if title.IdentificationNumber == "" {
			result = gin.H{
				"result":  "",
				"count":   id,
				"err":     "Title Not found",
				"success": 0,
			}
			c.JSON(http.StatusOK, result)
		} else {
			result = gin.H{
				"result":  title,
				"count":   id,
				"err":     "",
				"success": 1,
			}

			c.JSON(http.StatusOK, result)
		}

	})

	router.POST("/title", func(c *gin.Context) {
		//var buffer bytes.Buffer
		var title Title

		err = c.BindJSON(&title)
		if err != nil {
			fmt.Print(err.Error())
			//ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			c.AbortWithError(400, err)
			return
		}

		fmt.Println(title.IdentificationNumber)

		s := session.Copy()

		defer s.Close()

		col := s.DB(AuthDatabase).C(TitleCollection)

		err = col.Insert(&title)

		if err != nil {
			if mgo.IsDup(err) {
				c.JSON(http.StatusOK, gin.H{
					"result":  title,
					"count":   1,
					"err":     err.Error(),
					"success": 0,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"result":  title,
					"count":   1,
					"err":     "",
					"success": 1,
				})
			}
		}

		//_, err = stmt.Exec(openinglog.DayofweekID, openinglog.Date, openinglog.Drawopeningnumber, openinglog.Timegatesopen, openinglog.Timegatesdrop, openinglog.UserID, openinglog.BridgeID)

	})

	router.POST("/title/owner/update/:id", func(c *gin.Context) {
		//var buffer bytes.Buffer
		id := c.Param("id")
		var owner RegisteredOwner

		err = c.BindJSON(&owner)
		if err != nil {
			fmt.Print(err.Error())
			//ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			c.AbortWithError(400, err)
			return
		}

		fmt.Println(id)

		s := session.Copy()

		defer s.Close()

		col := s.DB(AuthDatabase).C(TitleCollection)

		docTitleToUpdate := bson.M{"identificationnumber": id}

		PushToRegisteredOwnerArray := bson.M{"$push": bson.M{"registeredowners": owner}}
		PushToOwnerArray := bson.M{"$set": bson.M{"owner": bson.M{"owner": owner}}}

		err = col.Update(docTitleToUpdate, PushToOwnerArray)

		if err != nil {
			fmt.Print(err.Error())
			c.AbortWithError(400, err)
			return
		}

		err = col.Update(docTitleToUpdate, PushToRegisteredOwnerArray)

		if err != nil {
			if mgo.IsDup(err) {
				c.JSON(http.StatusOK, gin.H{
					"result":  owner,
					"count":   1,
					"err":     err.Error(),
					"success": 0,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"result":  owner,
					"count":   1,
					"err":     "",
					"success": 1,
				})
			}
		}

		//_, err = stmt.Exec(openinglog.DayofweekID, openinglog.Date, openinglog.Drawopeningnumber, openinglog.Timegatesopen, openinglog.Timegatesdrop, openinglog.UserID, openinglog.BridgeID)

	})

	// GET all titles
	router.GET("/titles", func(c *gin.Context) {
		var (
		//title  Title
		//titles []Title
		)

		c.JSON(http.StatusOK, gin.H{
			//"result":  bridges,
			//"count":   len(bridges),
			"err":     "",
			"success": 1,
		})
	})

	// Add API handlers here
	router.Run(":5000")
}
