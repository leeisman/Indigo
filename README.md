# MQ

**mq 初始化**

```
mqInstance, err := mq.Init(ctx, topicID, credentialsFile, projectID, mq.InitPub())
if err != nil {
    log.Printf("failed to create client: %v", err)
}
```

**pub sample**

```
mqInstance.Publisher().
    Options(
        pub.SetErrorHook(func(err error, requestID string) {
            log.Printf("failed to publish requestID: %v ; message: %v", requestID, err)
        }),
    ).
    PublishJSON(pubMsg, requestID)
```

**pub options**
- SetErrorHook(func(err error,requestID string)) // 設定 error hoot，requestID 自定義，trace 方便。

**sub sample**

```
mqInstance.Subscriber().
    Options(
        sub.SyncMode(),
	    sub.SetErrorHook(func(err error) {
	        log.Printf("failed to pullMessageSync message: %v", err)
	).
    Subscribe(func(ctx context.Context, msg []byte) error {
	    log.Printf("Got Message: %s", string(msg))
		pubMsgJson, err := json.Marshal(pubMsg)
		if err != nil {
		    log.Printf(err.Error())
		}
    })
```
**sub options**
- SyncMode() // 同步。
- AsyncMode() // 非同步。
- SetMaxOutstandingMessages(i int) // 設定從topic 一次拉取的訊息數量。
- SetNumGoroutines(i int) // 設定async routing 數量。
- SetErrorHook(func(err error)) // 設定錯誤訊息hook。
