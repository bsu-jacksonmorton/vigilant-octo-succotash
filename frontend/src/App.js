import './App.css';
import {useState, useRef, useEffect} from 'react'
import Container from 'react-bootstrap/Container'
import 'bootstrap/dist/css/bootstrap.min.css'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import FormControl from 'react-bootstrap/FormControl'
import FormLabel from 'react-bootstrap/esm/FormLabel';
import FormGroup from 'react-bootstrap/esm/FormGroup';
import ListGroup from 'react-bootstrap/ListGroup'
import ListGroupItem from 'react-bootstrap/esm/ListGroupItem';

function App() {
  const [ messages, setMessages ] = useState([])
  const [ message, setMessage ] = useState("")
  const [ loggedIn, setLoggedIn ] = useState(false)
  const [ socketUrl, setSocketUrl ] = useState("-")
  const [ username, setUsername ] = useState("")
  const [webSocket, setWebSocket] = useState()


  function handleSendMessage(e){
    e.preventDefault()
    if(!loggedIn){
      alert("You are not currently connected to a server!")
      return
    }
    let messageToSend = {
      sender: username,
      body: message
    }
    setMessage("")
    webSocket.send(JSON.stringify(messageToSend))
  }
  function handleMessageUpdate(e){
    setMessage(e.target.value)
  }
  function handleLogin(e){
    e.preventDefault()
    let serverAddr = document.querySelector("#server-address-field").value
    let username = document.querySelector("#username-field").value
    if(!username.trim() || !serverAddr.trim()){
      alert("Invalid parameters provided.... try again!");
      return
    }
    // .... do websocket stuff?
    setSocketUrl(serverAddr)
    setUsername(username)
  }
  function handleLogout(e){
    e.preventDefault()
    setLoggedIn(false)
    setMessages([])
    setSocketUrl("")
    try{
      if(webSocket){
        webSocket.close()
      }
    }
    catch(e){
      console.log(e)
    }
  }
  useEffect(() => {
    if(socketUrl.indexOf("ws://") > -1){
      setWebSocket(new WebSocket(socketUrl))
    }
  }, [socketUrl])

  useEffect(() => {
    if(!webSocket) return
    webSocket.onopen = e => {
      console.log("connected to " + socketUrl)
      setLoggedIn(true)
    }
    webSocket.onclose = e => {
      console.log("disconnected from " + socketUrl)
      setLoggedIn(false)
    }
    webSocket.onmessage = e => {
      setMessages(messages.concat(JSON.parse(e.data)))
    }
  }, [webSocket, messages])
  return (
    <div className="App">
      <Container className="p-5">
        <Login username={username} serverAddress={socketUrl} loggedIn={loggedIn} handleLogin={handleLogin} handleLogout={handleLogout}></Login>
        <Messages messages={messages}></Messages>
        {loggedIn && <MessageInput handleMessageUpdate={handleMessageUpdate} handleSendMessage={handleSendMessage} message={message}></MessageInput>}
      </Container>
      
    </div>
  );
}

const Login = (props) => {
  const { loggedIn, serverAddress, username, handleLogin, handleLogout} = props
  if(!loggedIn){
    return(
      <Form onSubmit={handleLogin}>
        <FormGroup className="pb-2">
          <FormLabel>Server Address: </FormLabel>
          <FormControl as="input" id="server-address-field" defaultValue="ws://localhost:5005/ws" type="text"></FormControl>
        </FormGroup>
        <FormGroup className="pb-2">
          <FormLabel>Username: </FormLabel>
          <FormControl as="input" id="username-field" type="text"></FormControl>
       </FormGroup>
        <Button type="submit">connect</Button>
      </Form>
    )
  }
  else{
    return(
      <>
        <h3>connected to {serverAddress} as {username} </h3>
        <Button onClick={handleLogout}>disconnect</Button>
      </>
    )
  }
}
const Messages = (props) => {
  const { messages } = props
  return(
    <ListGroup className="mt-5">
      {messages.map(message => <ListGroupItem>{message.sender}: {message.body}</ListGroupItem>)}
    </ListGroup>
  )
}

const MessageInput = (props) => {
  const { handleMessageUpdate, handleSendMessage, message} = props
  return(
    <Form onSubmit={handleSendMessage}>
      <FormGroup className="mb-2">
        <FormLabel>say: </FormLabel>
        <FormControl as="input" value={message} onChange={handleMessageUpdate} type="text"></FormControl>
      </FormGroup>
      <Button type="submit">send</Button>
    </Form>
  )
}

export default App;
