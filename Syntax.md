```go
	content := []byte(utils.GetSystemPrompt(""))
	err = os.WriteFile("output.txt", content, 0644)
	if err != nil {
			panic(err)
	}

  println(utils.GetSystemPrompt(""))
  println(string(body))
```
