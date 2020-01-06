import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Header, Segment, Form, Message, Input, TextArea, Button, Table, Icon, Dropdown } from "semantic-ui-react";
//import CreateMessage from "./CreateMessage";

let endpoint = "http://127.0.0.1:8080";

// could make it a react funcitons
class ComposeForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      to: [],
      subject: "",
      body: "",
      scheduledAt: null,
      replySent: false,
      accountList: [],
      accountOptions: [],
      userAddedRecipients: [],
      _accountListTimer: 0
    }
  }

  componentDidMount() {
    this.getAccountList();
    this.enablePolling();
  }

  componentWillUnmount() {
    this.disablePolling()
  }

  enablePolling() {
    this.timer2 = setInterval(()=> this.getAccountList(), 3000)
    console.log("Enabling polling ",this.timer2)
    this.setState({_accountListTimer: this.timer2})
  }
  disablePolling() {
    console.log("Disabling polling ",this.state._accountListTimer)
    clearInterval(this.state._accountListTimer);
    this.timer2 = null
  }

  handleSubjectChange = (e) => {
    this.setState({
      subject: e.target.value
    });
  }

  handleBodyChange = (e) => {
    this.setState({
      body: e.target.value
    });
  }

  handleRecipientsChange = (e, {value}) => {
    console.log("HandleRecipients Change value e,b,c",e,value)
    this.setState({
      to: value
    });
  }

  handleAddRecipient = (e,{value}) => {
    console.log("Handle Add Recipient value", value)
    console.log("Handle Add Recipient State", this.state.userAddedRecipients)
    this.setState({
      userAddedRecipients: [...this.state.userAddedRecipients, value]
    },() => { this.updateAccountOptions() });
  }

  sentStatus = () => {
    if(this.state.replySent) {
      return   <Message success header='Sent' content='Your Message has been sent.'/>
    } else {
      return ""
    }
  }

  updateAccountOptions = () => {
    if(!this.props.IsLoggedIn) {
      return
    }

    // Load up options with just emails
    var newAccountOptions = this.state.accountList.map((account) => {
        var accOption = {};
        accOption.key = account.Email;
        accOption.text = account.Email;
        accOption.value = account.Email;
        return accOption;
      });
    // Load up Names
    newAccountOptions = [...newAccountOptions, ...this.state.accountList.map((account) => {
      var accOption = {};
      accOption.key = account.ID;
      accOption.text = account.FirstName+" "+account.LastName
      accOption.value = account.Email;
      return accOption;
    })];

    if(this.state.userAddedRecipients.length > 0) {
      newAccountOptions = [...newAccountOptions, ...this.state.userAddedRecipients.map((userVal) => {
        var accOption = {};
        accOption.key = userVal;
        accOption.text = userVal;
        accOption.value = userVal;
        return accOption;
      })];  
    }

    this.setState({ 
        accountOptions: newAccountOptions
      });
  }

  getAccountList = () => {
    let { IsLoggedIn } = this.props;
    if(!IsLoggedIn) {
      this.setState({accountList: []})
      return
    }

    console.log("===AccountList===")
    let apiStr = "/accountList"
    axios.get(endpoint + apiStr,
      {
        withCredentials: false
      } 
    ).then(res => {
      console.log(res);
      if(res.data) {
        this.setState({
          accountList: res.data
        }, () => { this.updateAccountOptions() });
      } else {
        this.setState({
          accountList: []
        });
      }
    },(error) => {
      console.log(error);
      this.disablePolling();
    });
  };


  sendMessage = () => {
    let a = this.props.GetAccountIDFn();
    let newSender = this.props.AccountEmail;
    let newRecipients = this.state.to;
    let newSubject = this.state.subject;

    let tostr = "";
    let apiStr = "/message";
    apiStr += "?"+a.name+"="+a.value;

    let data = JSON.stringify({ 
      ParentMid: 0,
      ScheduledAt: null,
      SenderEmail: newSender,
      Recipients: newRecipients,
      Subject: newSubject,
      Body: btoa(this.state.body)
    })
    console.log("Sending Message to", this.state.to)
    console.log("Sending Message",data)

    axios
      .post(endpoint + apiStr, data, 
        {
          headers: {
            "Content-Type": "application/json"
        }
      })
      .then(res => {
        console.log(res);
        this.setState({ 
          replySent: true, 
          body: ""
        });
      },(error)=>{
        console.log(error);
        this.setState({ 
          replySent: false
        });
      });    
  }

  createDisplay = () => {
    if(this.props.IsLoggedIn) {
      return (
      <Form success>
        <Table compact>
          <Table.Body>
            <Table.Row>
              <Table.Cell>From: </Table.Cell>            
              <Table.Cell>{this.props.AccountEmail}</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>To: </Table.Cell>
              <Table.Cell>
                <Dropdown
                  placeholder='Recipient, Recipient, ...'
                  fluid
                  multiple
                  search
                  selection
                  clearable
                  allowAdditions
                  value={this.state.to}
                  options={this.state.accountOptions}
                  onChange={this.handleRecipientsChange}
                  onAddItem={this.handleAddRecipient}
                />
                {/*<Input fluid id="to" placeholder="Recipient, Recipient, ..." rows={1} value={this.state.toRecipients} onChange={this.handleToChange}/>*/}
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Subject: </Table.Cell>
              <Table.Cell>
                <Input fluid id="subject" placeholder="Subject" rows={1} onChange={this.handleSubjectChange}/>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Body: </Table.Cell>
              <Table.Cell>
                <TextArea id="body" placeholder="Message" rows={3} onChange={this.handleBodyChange}/>
              </Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        {this.sentStatus()}
        <Button icon onClick={this.sendMessage}>
          <Icon name='send'/>Send</Button>
      </Form>
      );
    } else { 
      return (
        <></>
      );
    }
  }

  render() {
    return (
      <div>
        <Segment>
          <Header className="header" as="h3">{this.props.ComponentName}</Header>
          {this.createDisplay()}          
        </Segment>
      </div>
    );
  }
}



ComposeForm.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  AccountEmail: PropTypes.string.isRequired,

  ComponentName: PropTypes.string,
  FormatTimeFn: PropTypes.func
}

ComposeForm.defaultProps = {
  ComponentName: 'Compose Form',
  FormatTimeFn: identityFn
}

function identityFn(x) { return x; }

export default ComposeForm;
