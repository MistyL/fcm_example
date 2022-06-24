package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"fcm_example/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	cs := &http.Server{
		Addr:         ":8080",
		Handler:      clientServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	ss := &http.Server{
		Addr:         ":9090",
		Handler:      sendServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	g.Go(func() error {
		return cs.ListenAndServe()
	})
	g.Go(func() error {
		return ss.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

const server = "http://127.0.0.1:8000/v1"

func clientServer() http.Handler {
	r := gin.Default()

	r.Static("/css", "./html/css")
	r.Static("/js", "./html/js")
	r.Static("/img", "./html/img")

	r.LoadHTMLGlob("html/*.*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("username")
		c.Request.Method = "GET"
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"userId": name,
		})
	})

	r.POST("/registry", func(c *gin.Context) {
		var data models.RegistryUser
		c.BindJSON(&data)
		data.Platform = "web"
		body, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("=========json marshal error, %v\n", err)
			c.HTML(http.StatusBadGateway, "login.tmpl", gin.H{})
			return
		}
		req, err := http.NewRequest(http.MethodPost, server+"/registry", bytes.NewReader(body))
		if err != nil {
			fmt.Printf("==========new request error , %v\n", err)
			c.HTML(http.StatusBadGateway, "login.tmpl", gin.H{})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("=======server error, %v\n", err)
			c.HTML(http.StatusBadGateway, "login.tmpl", gin.H{})
			return
		} else {
			res, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			fmt.Printf("=======register user error, %v\n", string(res))
			if resp.StatusCode != 200 {
				c.HTML(http.StatusBadGateway, "login.tmpl", gin.H{})
				return
			} else {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{})
				return
			}
		}
	})

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/firebase-messaging-sw.js", func(c *gin.Context) {
		c.Header("Content-Type", "application/x-javascript")
		c.HTML(http.StatusOK, "firebase-messaging-sw.js", gin.H{})
	})

	return r
}

const adminServer = "http://127.0.0.1:8000/admin"

func sendServer() http.Handler {
	r := gin.Default()

	r.Static("/css", "./html/css")
	r.Static("/js", "./html/js")
	r.Static("/img", "./html/img")

	r.Use(Cors())

	r.LoadHTMLGlob("html/*.*")

	r.GET("/", func(c *gin.Context) {
		userResp, uErr := http.Get(adminServer + "/users")
		if uErr != nil {
			fmt.Printf("======get user list error , %v\n", uErr)
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}
		ubody, _ := ioutil.ReadAll(userResp.Body)
		defer userResp.Body.Close()
		if userResp.StatusCode != 200 {
			fmt.Printf("======request user list error, %v\n", string(ubody))
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}
		var u models.ListResp
		err := json.Unmarshal(ubody, &u)
		if err != nil {
			fmt.Printf("=======json unmarshal error, %v, %v\n", err, string(ubody))
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}

		topicResp, tErr := http.Get(adminServer + "/topics")
		if tErr != nil {
			fmt.Printf("=======get topic list error, %v\n", tErr)
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}
		tbody, _ := ioutil.ReadAll(topicResp.Body)
		defer topicResp.Body.Close()
		if topicResp.StatusCode != 200 {
			fmt.Printf("========reqeust topic list error, %v\n", string(tbody))
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}
		var t models.ListResp
		err = json.Unmarshal(tbody, &t)
		if err != nil {
			fmt.Printf("=======json unmarshal error, %v, %v\n", err, string(tbody))
			c.HTML(http.StatusBadGateway, "bar.tmpl", gin.H{})
			return
		}
		platform := []string{"web", "ios", "android"}
		c.HTML(http.StatusOK, "bar.tmpl", gin.H{
			"userList":  u.Users,
			"topicList": t.Topics,
			"platform":  platform,
		})
	})

	r.POST("/sendMessageToUser", func(c *gin.Context) {
		p := new(models.Paramters)
		c.BindJSON(&p)
		msg := models.GetDefaltMessage(p.UserId, "", "")

		resp, err := postServerWithBody(server+"/send/message", msg)
		if err != nil {
			c.String(http.StatusBadGateway, "send meessage error, %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			c.String(http.StatusBadGateway, "send message error, %s", string(body))
			return
		} else {
			c.HTML(http.StatusOK, "bar.tmpl", gin.H{})
			return
		}
	})

	r.POST("/sendMessageToTopic", func(c *gin.Context) {
		p := new(models.Paramters)
		c.BindJSON(&p)
		msg := models.GetDefaltMessage("", p.Topic, "")

		resp, err := postServerWithBody(server+"/send/message", msg)
		if err != nil {
			c.String(http.StatusBadGateway, "send message to topic error, %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			c.String(http.StatusBadGateway, "send message to topic error, %s", string(body))
			return
		} else {
			c.HTML(http.StatusOK, "bar.tmpl", gin.H{})
			return
		}
	})

	r.POST("/sendMessageToPlatform", func(c *gin.Context) {
		p := new(models.Paramters)
		c.BindJSON(&p)
		msg := models.GetDefaltMessage("", "", p.Platform)

		resp, err := postServerWithBody(server+"/send/message", msg)
		if err != nil {
			c.String(http.StatusBadGateway, "send message to platform error, %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			c.String(http.StatusBadGateway, "send message to platform error, %s", string(body))
			return
		} else {
			c.HTML(http.StatusOK, "bar.tmpl", gin.H{})
			return
		}
	})

	r.POST("/subscribeTopic", func(c *gin.Context) {
		p := new(models.SubscribeTopic)
		c.BindJSON(&p)
		body, err := json.Marshal(p)
		if err != nil {
			c.String(http.StatusBadGateway, "invalid paramter, %s", err.Error())
			return
		}
		resp, err := postServerWithBody(server+"/subscribe", body)
		if err != nil {
			c.String(http.StatusBadGateway, "server error, %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			c.String(http.StatusBadGateway, "response error, %s", string(body))
			return
		} else {
			c.HTML(http.StatusOK, "bar.tmpl", gin.H{})
			return
		}
	})

	r.POST("/unsubscribeTopic", func(c *gin.Context) {
		p := new(models.SubscribeTopic)
		c.BindJSON(&p)
		body, err := json.Marshal(p)
		if err != nil {
			c.String(http.StatusBadGateway, "invalid paramter, %s", err.Error())
			return
		}
		resp, err := postServerWithBody(server+"/unsubscribe", body)
		if err != nil {
			c.String(http.StatusBadGateway, "server error, %s", err.Error())
			return
		}
		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			c.String(http.StatusBadGateway, "response error, %s", string(body))
			return
		} else {
			c.HTML(http.StatusOK, "bar.tmpl", gin.H{})
			return
		}
	})

	return r
}

func postServerWithBody(url string, param []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(param))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("=========new request error, %v\n", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
