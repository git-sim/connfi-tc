import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Segment, Container, Grid, Card, Header, Label, Form, Input, Button, Menu, Icon } from "semantic-ui-react";


let endpoint = "http://127.0.0.1:8080";

class Folders extends Component {
  constructor(props) {
    super(props);
    // expect the following props
    //accountID: props.AccountIDValue,
    //accidName: props.AccountIDName,

    this.state = {
      folderid: 0,
      sort: 0,
      sortorder: -1,
      limit: 0,
      page: 0,
      task: "",
     folders: { 0: {name:"Inbox", count: 0, unviewed: 0},
                1: {name:"Archive", count: 0, unviewed: 0},
                2: {name:"Sent",count: 0, unviewed: 0},
                3: {name:"Scheduled",count: 0, unviewed:0 }
              },
      messages: [],
      _messageTimer: 0
    };

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
    this.timer = setInterval(()=> this.getFolderInfo(), 3000);
    console.log("Enabling polling ",this.timer)
    this.setState({_messageTimer: this.timer})
  }
  disablePolling() {
    console.log("Disabling polling ",this.state._messageTimer)
    clearInterval(this.state._messageTimer);
  }

  getFolderInfo = () => {
    let {limit, page, folderid } = this.state;
    let  AccountID = this.props.GetAccountIDFn();
    let  FolderID = this.props.FolderIDFn();
    let isLoggedIn = this.props.IsLoggedIn;

    if(!isLoggedIn) {
      return
    }

    let apiStr = "/folder";
    apiStr += "?"+AccountID.name+"="+AccountID.value;
    apiStr += "&"+FolderID.name +"="+folderid;
    axios.get(endpoint + apiStr,
      {
        params: {
          folderid,
          limit,
          page
        },
        withCredentials: false
      } 
    ).then(res => {
      console.log(res);
      if (res.data.Elems) {        
        this.setState({
          messages: res.data.Elems.map(msg => {
            let color = "yellow";

            if (msg.IsViewed) {
              color = "green";
            }
            return (
              <Card key={msg.Mid} color={color} fluid>
                <Card.Content>
                  <Card.Header textAlign="left">
                    <div style={{ wordWrap: "break-word" }}>{msg.M.M.Subject}</div>
                  </Card.Header>
                  <Card.Meta>{msg.M.SentAt}</Card.Meta>
                </Card.Content>
              </Card>
            );
          })
        });
      } else {
        this.setState({
          messages: []
        });
      }
    },(error) => {
      console.log(error);
      this.disablePolling();
    });
  };

  render() {
    return (
      <Grid rows={2}>
        <Segment>
          <Grid.Row>
            <Header className="header" as="h3">{this.props.ComponentName}</Header>
          </Grid.Row>
          <Grid.Row>                        
            <Menu vertical>
              <Menu.Item name="Inbox" 
                onClick={this.props.selectInbox} textalign="left" active={this.props.folderid===0}>
                Inbox
              </Menu.Item>
              <Menu.Item name="Archive" 
                onClick={this.props.selectArchive} textalign="left" active={this.props.folderid===1}>
                Archive
              </Menu.Item>
              <Menu.Item name="Sent" 
                onClick={this.props.selectSent} textalign="left" active={this.props.folderid===2}>
                Sent
              </Menu.Item>
              <Menu.Item name="Scheduled" 
                onClick={this.props.selectScheduled} textalign="left" active={this.props.folderid===3}>
                Scheduled
              </Menu.Item>
            </Menu>
          </Grid.Row>
        </Segment>
      </Grid>
    );
  }
}

Folders.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  GetFolderIDFn: PropTypes.func.isRequired,
  selectInbox: PropTypes.func.isRequired,
  selectArchive: PropTypes.func.isRequired,
  selectSent: PropTypes.func.isRequired,
  selectScheduled: PropTypes.func.isRequired
}

Folders.defaultProps = {
  ComponentName: 'Folders',
}

export default Folders;
