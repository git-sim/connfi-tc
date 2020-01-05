import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import {Segment, Grid, Pagination, Table, Header} from "semantic-ui-react";


let endpoint = "http://127.0.0.1:8080";

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
      _messageTimer: 0
    };

  }

  componentDidMount() {
    this.getMessages();
    this.enablePolling()
  }

  componentWillUnmount() {
    this.disablePolling()
  }

  enablePolling() {
    this.timer = setInterval(()=> this.getMessages(), 2000);
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
  
  getMessages = () => {
    let { sort, sortorder, limit, page } = this.state;
    let IsLoggedIn = this.props.IsLoggedIn;
    let FolderID = this.props.GetFolderIDFn();
    let AccountID = this.props.GetAccountIDFn();
    if(!IsLoggedIn) {
      this.setState({messages: []})
      return
    }

    console.log("===getMessages===")
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
      console.log(res);
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
            let timeStr = this.props.FormatTimeFn(msg.M.SentAt)
            return (
              <Table.Row key={msg.Mid} positive={!viewed} onClick={() => this.props.SetActiveMessageFn(msg)}>
                  <Table.Cell>{msg.M.M.SenderEmail}</Table.Cell>
                  <Table.Cell>{msg.M.M.Subject}</Table.Cell>
                  <Table.Cell>{timeStr}</Table.Cell>
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



Messages.propTypes = {
  IsLoggedIn: PropTypes.bool.isRequired,
  GetAccountIDFn: PropTypes.func.isRequired,
  GetFolderIDFn: PropTypes.func.isRequired,
  SetActiveMessageFn: PropTypes.func.isRequired
}

Messages.defaultProps = {
  ComponentName: 'Messages',
  FormatTimeFn: identityFn
}

function identityFn(x) { return x; }

export default Messages;
