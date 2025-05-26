package main

func main() {
	//localize := i18n.NewLocalize("etc/localize")
	//
	//languages := localize.GetSupportedLanguages()
	//fmt.Println(languages)
	//
	//text := localize.Get("text", "en")
	//fmt.Println(text)
	//
	//text = localize.Get("text", "zh")
	//fmt.Println(text)
	//
	//text = localize.Get("text", "zh-Hans")
	//fmt.Println(text)
	//
	//text = localize.Get("text", "ja")
	//fmt.Println(text)

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

	// aes := utils.NewAES("1234567812345678", "1234567812345678")
	// result, err := aes.Decrypt("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")
	// if err != nil {
	// 	println(err)
	// } else {
	// 	println(string(result))
	// 	// res, errE := aes.Decrypt(result)
	// 	// if errE != nil {
	// 	// 	println(errE)
	// 	// } else {
	// 	// 	println(string(res))
	// 	// }
	// }
}
