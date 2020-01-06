import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Segment, Grid, Button, Header, TextArea, Form, Message, Icon } from "semantic-ui-react";

// could make it a react functions
class CreateMessage extends Component {
  constructor(props) {
    super(props);

    this.state = {
      replyText: "",
      replySent: false
    }
  }

  handleReplyChange = (e) => {
    this.setState({
      replyText: e.target.value,
      replySent: false
    });
  }

  prepRecipients(origsender, origrecips, me) {
    let orecips = origrecips;
    for(let i = 0; i<orecips.length; i++){
      if(orecips[i]===me) {
        orecips[i]=origsender;
      }
    }
    return orecips;
  }

  sendResponse = (msg, respBody, isReplyAll) => {
    let origMsg = msg.M.M
    let a = this.props.GetAccountIDFn();
    let origSender = origMsg.SenderEmail
    let newSender = this.props.AccountEmail
    let newRecipients = this.prepRecipients(origSender,origMsg.Recipients,newSender)
    const replyPrefix = "Re:"
    let newSubject = ""

    if(!origMsg.Subject.startsWith(replyPrefix)) {
      newSubject = replyPrefix + origMsg.Subject;
    } else {
      newSubject = origMsg.Subject;
    }

    let apiStr = "/message";
    apiStr += "?"+a.name+"="+a.value;

    let data = JSON.stringify({ 
      ParentMid: msg.Mid,
      ScheduledAt: null,
      SenderEmail: newSender,
      Recipients: newRecipients,
      Subject: newSubject,
      Body: btoa(this.state.replyText)
    })

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
          replyText: ""
        });
      },(error)=>{
        console.log(error);
        this.setState({ 
          replySent: false
        });
      });    
  }

  sentSuccess = () => {
    if(this.state.replySent) {
      return   <Message success header='Sent' content='Your Reply has been sent.'/>
    } else {
      return ""
    }
  }


  displayNewMessageForm = () => {
    return (
      <>
        <Form success>
          <TextArea id="replyTextArea" placeholder="Reply" rows={3} value={this.state.replyText} onChange={this.handleReplyChange}/>
          {this.sentSuccess()}
          <Button icon onClick={() => 
            {this.sendResponse(this.props.ActiveMessage, this.state.replyText, false)}}>
            <Icon name='reply'/>Reply</Button>
          <Button icon onClick={() => 
            {this.sendResponse(this.props.ActiveMessage, this.state.replyText, true)}}>
            <Icon name='reply all'/>Reply All</Button>
        </Form>
      </>
    );
  }

  displayReplyMessageForm = () => {
    return (
      <>
        <Form success>
          <TextArea id="replyTextArea" placeholder="Reply" rows={3} value={this.state.replyText} onChange={this.handleReplyChange}/>
          {this.sentSuccess()}
          <Button icon onClick={() => 
            {this.sendResponse(this.props.ActiveMessage, this.state.replyText, false)}}>
            <Icon name='reply'/>Reply</Button>
          <Button icon onClick={() => 
            {this.sendResponse(this.props.ActiveMessage, this.state.replyText, true)}}>
            <Icon name='reply all'/>Reply All</Button>
        </Form>
      </>
    );
  }

  createDisplay = () => {
    if(this.props.IsLoggedIn) {
      if(this.props.ActiveMessage && 
        this.props.ActiveMessage.M) {
        // Active message passed in, show reply
        return this.displayReplyMessageForm()
      } else {
        // No active message, this is a new message
        return this.displayNewMessageForm();
      }
    } else { 
      return (
        <Grid rows={4}>
          <Grid.Row>
          </Grid.Row>
          <Segment>
          </Segment>
        </Grid>
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

// hack to get the public server working without adding react router
var endpoint = window.location.protocol+"//"+window.location.hostname+":8080"

CreateMessage.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  AccountEmail: PropTypes.string.isRequired,

  ActiveMessage: PropTypes.object,  //If ActiveMessage is null this is a new message
  FormatTimeFn: PropTypes.func
}

CreateMessage.defaultProps = {
  ComponentName: 'CreateMessage',
  FormatTimeFn: identityFn,
}

function identityFn(x) { return x; }

export default CreateMessage;
