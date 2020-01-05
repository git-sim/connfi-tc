import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Container, Segment, Grid, Button, Table, Header, TextArea, Form, Message, Icon } from "semantic-ui-react";

let endpoint = "http://127.0.0.1:8080";

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
    let a = this.props.GetAccountIDFn();
    let origSender = msg.M.M.SenderEmail
    let newSender = this.props.AccountEmail
    let newRecipients = this.prepRecipients(origSender,msg.M.M.Recipients,newSender)
    let apiStr = "/message";
    apiStr += "?"+a.name+"="+a.value;

    let data = JSON.stringify({ 
      ParentMid: msg.Mid,
      ScheduledAt: this.props.DefaultTimeString,
      SenderEmail: newSender,
      Recipients: newRecipients,
      Subject: "Re:"+msg.M.M.Subject,
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
      },(error)=>{
        console.log(error);
      }).then(data => {
        this.setState({ replySent: true, 
                        replyText: ""
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

  createDisplay = (showNew,showReply,showReplyAll) => {
    if(this.props.IsLoggedIn && 
      this.props.ActiveMessage && 
      this.props.ActiveMessage.M) {
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
          {this.createDisplay(/*showNew*/ false, /*showReply*/ true , /*showReplyAll*/ true)}
        </Segment>
      </div>
    );
  }
}



CreateMessage.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  ActiveMessage: PropTypes.object.isRequired,
  AccountEmail: PropTypes.string.isRequired,
  FormatTimeFn: PropTypes.func
}

CreateMessage.defaultProps = {
  ComponentName: 'CreateMessage',
  FormatTimeFn: identityFn,
  DefaulyTimeString: "2007-01-02T15:04:05Z07:00"

}

function identityFn(x) { return x; }

export default CreateMessage;
