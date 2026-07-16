package admin

import (
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type TestController struct {
}

func NewTestController() *TestController {
	return &TestController{}
}

func (c *TestController) Index(ctx *gin.Context) {
	//当前还没有 Operator，先生成测试 Token
	token, err := pkg.GenerateToken(
		10002,
		1,
		"nick",
	)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, gin.H{
		"token": token,
	})

	// client := skynet.New("127.0.0.1:8888")

	// if err := client.Connect(); err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	// system := adapter.NewSystemAdapter(client)

	// resp, err := system.Ping(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("msg      :", resp.Msg)
	// fmt.Println("timestamp:", resp.Timestamp)
}