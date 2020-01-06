import PropTypes from 'prop-types'
import React, { Component } from "react";
import axios from "axios";
import { Label, Menu } from "semantic-ui-react";

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
      folders: [ {key: 0, text:"Inbox",    value: 0, fn: this.props.selectInbox,     count: 0, unviewed: 0},
                 {key: 1, text:"Archive",  value: 1, fn: this.props.selectArchive,   count: 0, unviewed: 0},
                 {key: 2, text:"Sent",     value: 2, fn: this.props.selectSent,      count: 0, unviewed: 0},
                 {key: 3, text:"Scheduled",value: 3, fn: this.props.selectScheduled, count: 0, unviewed: 0}
      ],
      foldersLookup: {"Inbox":    0,
                      "Archive":  1,
                      "Sent":     2,
                      "Scheduled":3},
      messages: [],
      _messageTimer: 0
    };

    this.enablePolling = this.enablePolling.bind(this)
    this.disablePolling = this.disablePolling.bind(this)
    //this.getAccounts = this.getAccounts.bind(this)
  }

  componentDidMount() {
    this.enablePolling()
  }

  componentWillUnmount() {
    this.disablePolling()
  }

  enablePolling() {
    this.timer = setInterval(() => { this.getFolderCounts(); }, 900);
    console.log("Enabling polling ",this.timer)
    this.setState({_messageTimer: this.timer})
  }
  disablePolling() {
    console.log("Disabling polling ",this.state._messageTimer)
    clearInterval(this.state._messageTimer);
  }

  getFolderCounts = () => {
    let  AccountID = this.props.GetAccountIDFn();
    let  FolderID = this.props.GetFolderIDFn();
    let isLoggedIn = this.props.IsLoggedIn;

    if(!isLoggedIn) {
      return
    }

    let nFolders = this.state.folders.length;
    let apiStr = endpoint+"/folder";
    apiStr += "?"+AccountID.name+"="+AccountID.value;
    apiStr += "&limit=1&page=0&sort=0&sortorder=-1";
    apiStr += "&"+FolderID.name +"=";

    // Get the info (message counts) for the folders
    let fGets = [];
    for (let i = 0; i < nFolders; i++) {
      fGets.push(axiosGet(apiStr+i));
    }

    Promise.all(fGets).then(fRes => {
      this.setState({
        folders: this.state.folders.map( (folder) => {
          folder.count =  fRes[folder.key].data.NumTotal;
          folder.unviewed = fRes[folder.key].data.NumUnviewed;
          return folder;      
        })
      });
    });
  }    

  displayFolderMenu = () => {
    if(this.props.IsLoggedIn) {
      return ( this.state.folders.map(folder => {
        return (
          <Menu.Item 
            key={folder.key} 
            name={folder.text}
            active={this.state.folderid===folder.value}
            onClick={this.selectFolder}>
              {folder.text}
              { folder.key===0 && 
                <Label position="right" 
                  color={(() => 
                    {if(this.state.folders[folder.key].unviewed>0) 
                      { return 'teal'} 
                      else {return 'grey'}
                    })()}>
                {folder.unviewed}
                </Label>
              }
              <Label position="right" color="grey">
                {folder.count}
              </Label>
            </Menu.Item>
          );
        })
      );
    } else {
      return ( this.state.folders.map(folder => {
        return(           
          <Menu.Item 
            key={folder.key} 
            name={folder.text}
            active={this.state.folderid===folder.value}>
              {folder.text}
            { (folder.key===0) && <Label position="right" color="grey">0</Label> }
            <Label position="right" color="grey">0</Label>
          </Menu.Item>
        );
        })
      );
    }
  }

  selectFolder = (e,data) => {
    let idx = this.state.foldersLookup[data.name]
    this.setState({
      folderid: idx
    });
    this.state.folders[idx].fn();
  }


  render() {
    return (
      <Menu
        defaultActiveIndex={0}>
        {this.displayFolderMenu()}
      </Menu>
    );
  }
}

function axiosGet(apiStr) {
  return axios.get(apiStr,{withCredential: false});
}

// hack to get the public server working without adding react router
var endpoint = window.location.protocol+"//"+window.location.hostname+":8080"

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
