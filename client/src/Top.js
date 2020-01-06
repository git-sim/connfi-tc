import React, { Component } from "react";
import axios from "axios";
import { Segment, Grid, Container, Header, Form, Input, Button, Menu } from "semantic-ui-react";
import ComposeForm from "./ComposeForm";
import Folders from "./Folders";
import Messages from "./Messages";
import MessageView from "./MessageView";

const AccountIDName = "accid"
const FolderIDName = "folderid"
const MessageIDName = "msgid"
const ViewedName = "viewed"
//const StarredName = "starred"

class Top extends Component {
  constructor(props) {
    super(props);

    this.state = {
	  endpoint: endpoint,
      email: "",
      isComposing: false,
      folderid: 0,
      sort: 0,
      sortorder: -1,
      limit: 10,
      page: 0,
      Account: {
        ID:"",
        Email:"",
        FirstName:"",
        LastName:""
      },
      isLoggedIn: false,
      task: "",
      items: [],
      /*folders: [{name:"Inbox",val:0},
                {name:"Archive",val:1},
                {name:"Sent",val:2},
                {name:"Scheduled",val:3}
              ],*/
      messages: [],
      activeMessage: {},
      _messageTimer: 0
    };

    this.onChange = this.onChange.bind(this)
    this.onLogInOut = this.onLogInOut.bind(this)
    this.enablePolling = this.enablePolling.bind(this)
    this.disablePolling = this.disablePolling.bind(this)
    //this.getAccounts = this.getAccounts.bind(this)
  }

  componentDidMount() {
    //this.getAccounts();
    //this.enablePolling()
  }

  componentWillUnmount() {
    this.disablePolling()
  }

  enablePolling() {
    this.timer = setInterval(()=> this.setState({}), 900);
    console.log("Enabling polling ",this.timer)
    this.setState({_messageTimer: this.timer})
  }
  disablePolling() {
    console.log("Disabling polling ",this.state._messageTimer)
    clearInterval(this.state._messageTimer);
  }

  onChange = event => {
    console.log(event.target.name, event.target.value);
    this.setState({
      email: event.target.value
    });
  };

  onLogInOut = () => {
    let { email  } = this.state;
    let { isLoggedIn } = this.state;
    if(isLoggedIn === false) {
      if (email) {
        console.log("email onsubmit", this.state.email);
        axios
          .post(
            endpoint + "/login", 
            null,
            {
              params: {
                email
              },
              headers: {
                "Content-Type": "application/x-www-form-urlencoded"
              },
              withCredentials: false
            }
          )
          .then(res => {         
              console.log("OnSubmit response ",res);
              //cookie.save("session-cookie",res.headers["session-cookie"]);
              this.setState({
                Account: res.data,
                isLoggedIn: true
              });
              this.enablePolling();
            }, 
            (error) => { console.log("OnSubmit error", error); }
          );
      }
    } else {
      axios
        .post( endpoint + "/logout")
        .then(res => {
          console.log(res)
        }, (error) => {
          console.log(error)
        });
      this.setState({
        Account: {},
        isLoggedIn: false,
        email: "",
        messages: []
      });
      this.disablePolling();
    }
  };

  setFolderid = (idx) => {
    this.setState({
      folderid: idx,
      isComposing: false
    })
  }

  composeClick = () => {
    this.setState({ isComposing: true });
  }


  // the GetxxxIDObj are to isolate the lower level components from 
  // details of the actual names of the url params (ie 'accid','folderid','msgid').
  GetFolderIDObj = () => {
    var id = { name: FolderIDName, value: this.state.folderid };
    return id;
  }

  GetAccIDObj = () => {
    var id = { name: AccountIDName, value: ""};
    if(this.state.Account) {
      id.value = this.state.Account.ID;
    }
    return id;
  }

  GetMessageIDObj = () => {
    var id = { name: MessageIDName, value: ""};
    if(this.state.activeMessage) {
      id.value = this.state.activeMessage.Mid;
    }
    return id;
  }

  GetViewedObj = () => {
    var o = { name: ViewedName, value: 0};
    return o
  }

  GetAccIDStr = () => {
    let accid = ""
    if(this.state.Account) {
      accid = this.state.Account.ID
    }
    return accid 
  }

  markAsViewed = (msg, val) => {
    let a = this.GetAccIDObj();
    let m = this.GetMessageIDObj();
    let v = this.GetViewedObj();
    let apiStr = "/message";
    apiStr += "?"+a.name+"="+a.value;
    apiStr += "&"+m.name+"="+msg.Mid.toString(16);
    apiStr += "&"+v.name;
    if(val) {
      apiStr += "=1";
    }
    else {
      apiStr += "=0";
    }

    axios
      .put(endpoint + apiStr, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        console.log(res);
      });
    
  }

  setActiveMessage = (msg) => {
    this.setState({ activeMessage: msg });
    this.markAsViewed(msg,true);
  }

  renderComposeOrMessages = () => {
    if(this.state.isComposing) {
      return (
        <Container fluid>
          <ComposeForm
            ComponentName="New Message"
            IsLoggedIn={this.state.isLoggedIn} 
            GetAccountIDFn={this.GetAccIDObj}
            AccountEmail={this.state.email}
            />
        </Container>            
      );
    } else {
      return (
        <Segment>
          <Messages 
          ComponentName="Messages"
          IsLoggedIn={this.state.isLoggedIn}
          GetAccountIDFn={this.GetAccIDObj}
          GetFolderIDFn={this.GetFolderIDObj}
          FormatTimeFn={formatGoTime}
          SetActiveMessageFn={(msg) => {this.setActiveMessage(msg)}}
          ActiveMessage={this.state.activeMessage}
          />
        </Segment>
      );
    }

  }

  render() {
    return (
      <div>
        <Grid rows={4}>
          <Grid.Row height={5}>
            <Header className="header" as="h2">TC Messaging</Header>
          </Grid.Row>
          <Grid.Row height={5}>
            <Menu fluid>
              <Menu.Item>
                Image
              </Menu.Item>
              <Menu.Item>
                Welcome {this.state.Account.Firstname} {this.state.Account.Firstname} {this.state.Account.Email}
              </Menu.Item>
              <Menu.Item position='right'>
                <Form onSubmit={this.onLogInOut}>
                    <Input
                      type="text"
                      name="email"
                      onChange={this.onChange}
                      value={this.state.email}
                      placeholder="Email Address"
                    />
                    <Button >{this.state.isLoggedIn ? "Logout" : "Login"}</Button>
                </Form>
              </Menu.Item>
            </Menu>
          </Grid.Row>
          <Grid.Row>
            <Menu fluid>
              <Folders 
                  IsLoggedIn={this.state.isLoggedIn} 
                  GetAccountIDFn={this.GetAccIDObj}
                  GetFolderIDFn={this.GetFolderIDObj}
                  selectInbox={() => {this.setFolderid(0)} }
                  selectArchive={() => {this.setFolderid(1)} }
                  selectSent={() => {this.setFolderid(2)} }
                  selectScheduled={() => {this.setFolderid(3)} }/>
              <Menu floated="right">
                <Menu.Item>
                  <Button onClick={this.composeClick}>Compose</Button>
                </Menu.Item>
              </Menu>
            </Menu>
            {this.renderComposeOrMessages()}
          </Grid.Row>
          <Grid.Row columns={1} height={50}>
            <Container fluid>
              <MessageView 
                ComponentName="View"
                IsLoggedIn={this.state.isLoggedIn}
                GetAccountIDFn={this.GetAccIDObj}
                ActiveMessage={this.state.activeMessage}
                AccountEmail={this.state.email}
                FormatTimeFn={formatGoTime}/>
              </Container>
          </Grid.Row>
        </Grid>
      </div>
    );
  }
}

function formatGoTime(instr) {
  // Todo use a real format routine that takes the locale formatting into account
  // Example GO time.Time as a json string
  // 2020-01-04T10:35:58.8690175-07:00
  // 0123456789012345678901234567890123
  // 0         1         2         3
  //const year = instr.substring(0,4);
  //const month = instr.substring(5,7);
  //const day = instr.substring(8,10);
  //const time = instr.substring(11,19);
  //return time+" "+month+" "+day+" "+year;
  const dt = new Date(instr);
  return dt.toLocaleString();

}

// hack to get the public server working without adding react router
var endpoint = window.location.protocol+"//"+window.location.hostname+":8080"

export default Top;
