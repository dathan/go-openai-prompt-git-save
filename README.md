## Purpose

[https://dathan.github.io/blog/posts/chat-gpt-4-save-prompts-to-git/](https://dathan.github.io/blog/posts/chat-gpt-4-save-prompts-to-git/)

## Requirements 

* OPENAI_API_KEY
* GITHUB_TOKEN

Environment variables must be defined otherwise the command line application will panic

## Features
* Makefile to build consistently in a local environment and remote environment
* Dockerfile for a generic image to build for 
* Go Mod (which you should to your project path change)
* VS Code environment
* Generic docker push

## TODO
* Brew generic install [DONE]
* GITHUB Actions build and push to dockerhub [DONE]
* Production Builds with git tag

## Installing via brew
* `brew install --verbose --build-from-source brew/Formula/go-openai-prompt-git-save.rb`
