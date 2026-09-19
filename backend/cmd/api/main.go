package main

import("log";"os";"github.com/sagarbagwe/jobapply-ai/backend/internal/ai";"github.com/sagarbagwe/jobapply-ai/backend/internal/httpapi")
func main(){addr:=env("HTTP_ADDR",":8080");provider:=ai.NewGemini(os.Getenv("GEMINI_API_KEY"),env("GEMINI_MODEL","gemini-2.5-flash"));log.Printf("JobApply AI API listening on %s",addr);if err:=httpapi.Router(provider).Run(addr);err!=nil{log.Fatal(err)}}
func env(k,f string)string{if v:=os.Getenv(k);v!=""{return v};return f}
