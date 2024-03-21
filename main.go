package main

import "github.com/thk-im/thk-im-base-server/utils"

func main() {

	// config := &conf.Config{}
	// err := conf.Load("etc/server.yaml", config)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// srvContext := &server.Context{}
	// srvContext.Init(config)
	//
	// for i := 0; i < 1000; i++ {
	// 	srvContext.Logger().Info("logger message:", i)
	// }
	//
	// time.Sleep(time.Hour)

	// keys := make([]string, 0)
	// for i := 0; i < 10000; i++ {
	// 	keys = append(keys, utils.GetRandomString(12))
	// }
	// resp, errGet := utils.BatchGetString(srvContext.RedisCache(), keys)
	// if errGet != nil {
	// 	fmt.Println(errGet)
	// } else {
	// 	fmt.Println("len:", len(resp))
	// }
	//
	// errSet := srvContext.RedisCache().Set(context.Background(), "22", 11, time.Hour).Err()
	// if errSet != nil {
	// 	fmt.Println(errSet)
	// }
	//
	// ok, errDel := utils.DelKeyByValue(srvContext.RedisCache(), "22", 11)
	// if errDel != nil {
	// 	fmt.Println(errDel)
	// } else {
	// 	fmt.Println("ok: ", ok)
	// }

	// objectPath := "etc/test.png"
	// now := time.Now().UnixMilli()
	// urlPath, errUpload := srvContext.ObjectStorage().UploadObject(fmt.Sprintf("%d-test.png", now), objectPath)
	// if errUpload != nil {
	// 	fmt.Println(errUpload)
	// } else {
	// 	fmt.Println(urlPath)
	// }

	// token, errToken := utils.GenerateUserToken(1, "12313", "11111")
	// if errToken != nil {
	// 	fmt.Println(errToken)
	// } else {
	// 	fmt.Println(token)
	// }
	//
	// id, errId := utils.CheckUserToken(token, srvContext.Config().Cipher)
	// if errId != nil {
	// 	fmt.Println(errId)
	// } else {
	// 	fmt.Println(id)
	// }

	aes := utils.NewAES("1231311212313112", "1231311212313112")
	result, err := aes.Encrypt([]byte("123123"))
	if err != nil {
		println(err)
	} else {
		println(result)
		res, errE := aes.Decrypt(result)
		if errE != nil {
			println(errE)
		} else {
			println(string(res))
		}
	}
}
