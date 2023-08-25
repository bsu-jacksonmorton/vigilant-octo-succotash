# Chat Server

* Author: Jackson Morton

## Overview
This project is a simple little chat server that currently acts as a medium for me to learn some new languages/frameworks. Feel free to check it out or use as starting code / inspiration for your own project :)

## Compiling and Using
### Server 
To compile and run the server program, use the following commands from the project root directory run:

    $ go run backend/main.go

After running that commmand you should see the following statement printed to std out:

    Server has started on port 5005!

### Client
To get the client up and running, use the following commands from the project root directory:

    $ cd frontend
    $ npm i
    $ npm start

This should start the client process and you should be able to access on localhost:3000 in your browser.

## Discussion
This project is a simple chat application that I am working on as a means of learning some new languages/frameworks (in particular Go and ReactJS).

The backend server is written in Go. I was interested in picking up Go as a language due to the great things I had heard about its multi-threading support and I figured a chat server would be a great way to test those features. Being able to spin up a thread by simply using the keyword "go" is very satisfying especially when you compare it to the boilerplate you would need in a language like Java. I prefer the way that error handling is managed in Go compared to other languages where exceptions are used.

The frontend client is being built with react + bootstrap. React definitely has a learning curve but I like what I am seeing so far and I can already see why it would be the weapon of choice for building out the frontend in industry. Components make things nice and modular and useState solves most of the headache I run into when it comes to managing the DOM. Right now the frontend is.... ugly but it is definitely responsive / reactive!



