import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import {Segment, Grid, Pagination, Table, Icon, Header} from "semantic-ui-react";

class Messages extends Component {
  constructor(props) {
    super(props);

    this.state = {
      sort: 0,
      sortorder: -1,
      limit: 10,
      page: 0,
      task: "",
      nMsgsInFolder: 0,
      nUnviewedInFolder: 0,
      messages: [],
      selectedMessage: null,
      _messageTimer: 0,
    };

  }

  componentDidMount() {
    //this.getMessages();
    this.enablePolling();
  }

  componentWillUnmount() {
    this.disablePolling();
  }

  enablePolling() {
    this.timer = setInterval(()=> this.getMessages(), 900);
    console.log("Enabling polling ",this.timer)
    this.setState({_messageTimer: this.timer})
  }
  disablePolling() {
    console.log("Disabling polling ",this.state._messageTimer)
    clearInterval(this.state._messageTimer);
  }

  onPageChange = (event, data) => {
    let ap = data["activePage"]
    let newpage = 0
    console.log("onPageChange ", ap, data)
    if(ap > 0) {
      newpage = ap-1 
    } else {
      newpage = 0
    };

    this.setState( { page: newpage });
    this.getMessages();
    console.log("called getmessages")
  }

  getNumPages = () => {
    let messagesPerPage = this.state.limit;
    if(this.state.messages && messagesPerPage > 0) {
      let answer = this.state.nMsgsInFolder/messagesPerPage;
      if(answer < 1) {
        return 1;
      }
      return answer;
    }
    return 1;
  }
  
  getIsActiveMsg = (mid) => {
    let isLoggedIn = this.props.IsLoggedIn;
    let act = this.props.ActiveMessage;
    if(isLoggedIn && act && parseInt(act.Mid,10) === parseInt(mid,10)) {
      return true;
    }
    return false;
  }

  getIconName = (isViewed) => {
    if(isViewed) {
      return "envelope open outline"
    } else {
      return "envelope"
    }
  }

  getIconColor = (isViewed) => {
    if(isViewed) {
      return "grey"
    } else {
      return "teal"
    }

  }

  getMessages = () => {
    let { sort, sortorder, limit, page } = this.state;
    let IsLoggedIn = this.props.IsLoggedIn;
    let FolderID = this.props.GetFolderIDFn();
    let AccountID = this.props.GetAccountIDFn();
    if(!IsLoggedIn) {
      this.setState({messages: []})
      return
    }

    //console.log("===getMessages===")
    let apiStr = "/folder"
    apiStr += "?"+AccountID.name+"="+AccountID.value
    apiStr += "&"+FolderID.name+"="+FolderID.value
    axios.get(endpoint + apiStr,
      {
        params: {
          sort,
          sortorder,
          limit,
          page
        },
        withCredentials: false
      } 
    ).then(res => {
      //console.log(res);
      if(res.data) {
        this.setState({
          nMsgsInFolder: res.data.NumTotal,
          nUnviewedInFolder: res.data.NumUnviewed
        });
      }
      if (res.data.Elems) {        
        this.setState({
          messages: res.data.Elems.map(msg => {
            let viewed = msg.IsViewed;
            let timeStr = this.props.FormatTimeFn(msg.M.SentAt);
            let base = msg.M.M;
            return (
              <Table.Row key={msg.Mid} 
              active={this.getIsActiveMsg(msg.Mid)} 
              positive={!viewed}
              onClick={() => this.props.SetActiveMessageFn(msg)}>
                  <Table.Cell>{base.SenderEmail}</Table.Cell>
                  <Table.Cell>{base.Subject}</Table.Cell>
                  <Table.Cell>{timeStr} 
                      <Icon 
                        name={this.getIconName(viewed)} 
                        color={this.getIconColor(viewed)}/>
                  </Table.Cell>
              </Table.Row>
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
      <div>
        <Grid>
          <Segment>
            <Header className="header" as="h3">{this.props.ComponentName}</Header>                        
            <Table fixed selectable compact>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>From</Table.HeaderCell>
                  <Table.HeaderCell>Subject</Table.HeaderCell>
                  <Table.HeaderCell>Date</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>{this.state.messages}</Table.Body>
            </Table>
            <Pagination pointing secondary
              disabled={false} 
              totalPages={this.getNumPages()} 
              onPageChange={this.onPageChange}
              />
          </Segment>
        </Grid>
      </div>
    );
  }
}

// hack to get the public server working without adding react router
var endpoint = window.location.protocol+"//"+window.location.hostname+":8080"


Messages.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  GetFolderIDFn: PropTypes.func.isRequired,
  SetActiveMessageFn: PropTypes.func.isRequired,
  ActiveMessage: PropTypes.object.isRequired
}

Messages.defaultProps = {
  ComponentName: 'Messages',
  FormatTimeFn: identityFn
}

function identityFn(x) { return x; }

export default Messages;
