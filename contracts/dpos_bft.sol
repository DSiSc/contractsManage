pragma solidity >=0.4.24 <0.5.0;

// DPOS contract interface
interface IDposContract {
    // set the initial member list
    function initMembers(address[] members) external;
    // init the member url
    function initMemberUrl(address addr, string url) external;
    // init consensus timeout config
    function initConsensusTimer(uint64 blockDelay, uint64 idleTimeOut, uint64 commitTimeOut, uint64 viewChangeTimeOut) external;
    // vote to add a new member
    function voteAddMember(address addr) external;
    // vote to remove an existing member
    function voteRemoveMember(address addr) external;
    // update member url
    function updateMemberUrl(string url) external;
    // get current member list
    function getMembers() external view returns (address[]);
    // get current member list
    function getMemberUrl(address addr) external view returns (string);
    // get consensus timer info
    function getConsensusTimer() external view returns (uint64, uint64, uint64, uint64);
}

contract SimpleDposContract is IDposContract {
    struct consensusTimer {
        uint64	blockDelay;
        uint64	idleTimeOut;
        uint64	commitTimeOut;
        uint64	viewChangeTimeOut;
    }
    consensusTimer consensus;
    address[] validMemberAddrs;
    mapping(address => string) validMemberInfo;
    mapping(address => address[]) voteAddMembers;
    mapping(address => address[]) voteRemoverMembers;

    function initMembers(address[] members) external {
        require(arraryLen(validMemberAddrs) == 0);
        validMemberAddrs = members;
    }

    function initMemberUrl(address addr, string url) external {
        require(existing(validMemberAddrs, addr));
        require(!existingUrl(addr));
        validMemberInfo[addr] = url;
    }

    function initConsensusTimer(uint64 blockDelay, uint64 idleTimeOut, uint64 commitTimeOut, uint64 viewChangeTimeOut) external {
        require(idleTimeOut >= 1000);
        require(commitTimeOut >= 1000);
        require(viewChangeTimeOut >= 5000);
        require(blockDelay + 1000 <= idleTimeOut);
        require(consensus.idleTimeOut == 0);
        require(consensus.commitTimeOut == 0);
        require(consensus.viewChangeTimeOut == 0);
        consensus.blockDelay = blockDelay;
        consensus.idleTimeOut = idleTimeOut;
        consensus.commitTimeOut = commitTimeOut;
        consensus.viewChangeTimeOut = viewChangeTimeOut;
    }

    function voteAddMember(address addr) external {
        require(isValidVoter(msg.sender));
        if (existing(validMemberAddrs, addr)) {
            return;
        }
        if (existing(voteAddMembers[addr], msg.sender)) {
            return;
        }

        addElement(voteAddMembers[addr], msg.sender);
        if (voteAddMembers[addr].length > arraryLen(validMemberAddrs) / 2) {
            addElement(validMemberAddrs, addr);
            voteAddMembers[addr] = new address[](0);
        }
    }

    function voteRemoveMember(address addr) external {
        require(isValidVoter(msg.sender));
        if (!existing(validMemberAddrs, addr)) {
            return;
        }
        if (existing(voteRemoverMembers[addr], msg.sender)) {
            return;
        }

        addElement(voteRemoverMembers[addr], msg.sender);

        if (voteRemoverMembers[addr].length > arraryLen(validMemberAddrs) / 2) {
            removeValidMember(addr);
            voteRemoverMembers[addr] = new address[](0);
        }
    }

    function updateMemberUrl(string url) external {
        require(existing(validMemberAddrs, msg.sender));
        validMemberInfo[msg.sender] = url;
    }

    function getMembers() external view returns (address[]) {
        address[] memory validAddrs = new address[](validMemberAddrs.length);
        for (uint i = 0; i<validMemberAddrs.length; i++){
            if (existingUrl(validMemberAddrs[i])) {
                validAddrs[i] = validMemberAddrs[i];
            }
        }
        return validAddrs;
    }

    function getMemberUrl(address addr) external view returns (string) {
        return validMemberInfo[addr];
    }

    // get consensus timer info
    function getConsensusTimer() external view returns (uint64, uint64, uint64, uint64) {
        return (consensus.blockDelay, consensus.idleTimeOut, consensus.commitTimeOut, consensus.viewChangeTimeOut);
    }

    function removeValidMember(address addr) internal {
        uint index = validMemberAddrs.length;
        for (uint i = 0; i<validMemberAddrs.length; i++){
            if (validMemberAddrs[i] == addr) {
                index = i;
                break;
            }
        }
        delete validMemberAddrs[index];
        validMemberInfo[addr] = "";
    }

    function addElement(address[] storage addrList, address addr) internal {
        for (uint i = 0; i<addrList.length; i++){
            if (addrList[i] == address(0)) {
                addrList[i] = addr;
                return;
            }
        }
        addrList.push(addr);
    }

    function arraryLen(address[] storage addrList) internal view returns(uint) {
        uint length = 0;
        for (uint i = 0; i<addrList.length; i++){
            if (addrList[i] != address(0)) {
                length++;
            }
        }
        return length;
    }

    function isValidVoter(address addr) internal view returns(bool) {
        return existing(validMemberAddrs, addr);
    }

    function existing(address[] storage addrList, address addr) internal view returns (bool) {
        for (uint i = 0; i<addrList.length; i++){
            if (addrList[i] == addr) {
                return true;
            }
        }
        return false;
    }

    function existingUrl(address addr) internal view returns (bool) {
        bytes memory tempStrBytes = bytes(validMemberInfo[addr]);
        if (tempStrBytes.length <=0 ) {
            return false;
        }
        return true;
    }
}