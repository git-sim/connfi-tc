import PropTypes from 'prop-types'
import React, { Component } from "react";
import { Header, Segment, Form, Message, Input, TextArea, Menu, Button, Table, Icon } from "semantic-ui-react";
import CreateMessage from "./CreateMessage";

// could make it a react funcitons
class ComposeForm extends Component {
  constructor(props) {
    super(props);

    this.state = {
      replyText: "",
      replySent: false
    }
  }

  handleToChange = (e) => {
    this.setState({
      replyText: e
    });
  }

  sentSuccess = () => {
    if(this.state.replySent) {
      return   <Message success header='Sent' content='Your Reply has been sent.'/>
    } else {
      return ""
    }
  }

  sendMessage = () => {
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
                <Input fluid id="to" placeholder="Recipient, Recipient, ..." rows={1} value={this.state.toRecipients} onChange={this.handleToChange}/>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Subject: </Table.Cell>
              <Table.Cell>
                <Input fluid id="subject" placeholder="Subject" rows={1}/>
              </Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Body: </Table.Cell>
              <Table.Cell>
                <TextArea id="body" placeholder="Message" rows={3}/>
              </Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        {this.sentSuccess()}
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
