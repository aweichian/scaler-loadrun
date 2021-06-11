package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func Health(c *gin.Context) {
	c.AsciiJSON(http.StatusOK, "ok")
}

type Response struct {
	Type     string
	Duration string
}

func Task(c *gin.Context) {
	method := c.Request.Method
	var result float64
	var response Response
	if method == "GET" {
		typ := c.Query("type")
		number := c.Query("number")
		num, _ := strconv.Atoi(number)
		start := time.Now()
		switch typ {
		case "cpu":
			result = Cpu(float64(num))
		case "memory":
			result = Memory(num)
		case "all":
			All(num)
		default:

		}
		fmt.Printf("result: %v\n", result)
		diff := time.Now().Sub(start)
		fmt.Printf("time used: %v\n", diff)
		response = Response{
			Type:     typ,
			Duration: diff.String(),
		}
	}

	c.JSON(200, response)

}

func Cpu(num float64) (sum float64) {
	var i float64
	for i = 0; i < num; i++ {
		sum += i
	}
	return sum
}

type Info struct {
	Name   string
	Age    int
	Gender string
	Salary float64
}

func Memory(num int) (count float64) {
	var infos []Info

	for i := 0; i < num; i++ {
		infos = append(infos, Info{
			Name:   fmt.Sprintf("awei%d", i),
			Age:    i,
			Gender: "female",
			Salary: float64(i),
		})
	}
	return float64(len(infos))
}

func All(num int) {
	var infos []Info
	var sum float64
	for i := 0; i < num; i++ {
		infos = append(infos, Info{
			Name:   fmt.Sprintf("awei%d", i),
			Age:    i,
			Gender: "female",
			Salary: float64(i),
		})
		sum += float64(i)
	}
}
