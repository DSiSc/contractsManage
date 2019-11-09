pragma solidity >=0.4.24 <0.5.0;

// MSP contract interface
interface IMspContract {
    // set the initial manager list
    function initManagers(address[] members) external;
    // vote to add a new manager
    function voteAddManager(address addr) external;
    // vote to remove an existing manager
    function voteRemoveManager(address addr) external;
    // authorize new member
    function authorizeMember(address addr) external;
    // revoke authorization
    function revokeAuthorization(address addr) external;
    // check whether member is authorized
    function isAuthorized(address addr) external view returns (bool);
}

contract SimpleMspContract is IMspContract {
    address[] validManagers;
    mapping(address => address[]) voteAddManagers;
    mapping(address => address[]) voteRemoverManagers;
    mapping(address => bool)      authorizedMembers;

    function initManagers(address[] members) external {
        require(arraryLen(validManagers) == 0);
        validManagers = members;
    }

    function voteAddManager(address addr) external {
        require(isValidVoter(msg.sender));
        if (existing(validManagers, addr)) {
            return;
        }
        if (existing(voteAddManagers[addr], msg.sender)) {
            return;
        }

        addElement(voteAddManagers[addr], msg.sender);
        if (voteAddManagers[addr].length > arraryLen(validManagers) / 2) {
            addElement(validManagers, addr);
            voteAddManagers[addr] = new address[](0);
        }
    }

    function voteRemoveManager(address addr) external {
        require(isValidVoter(msg.sender));
        if (!existing(validManagers, addr)) {
            return;
        }
        if (existing(voteRemoverManagers[addr], msg.sender)) {
            return;
        }

        addElement(voteRemoverManagers[addr], msg.sender);

        if (voteRemoverManagers[addr].length > arraryLen(validManagers) / 2) {
            removeValidMember(addr);
            voteRemoverManagers[addr] = new address[](0);
        }
    }

    function authorizeMember(address addr) external {
        require(isValidVoter(msg.sender));
        authorizedMembers[addr] = true;
    }

    function revokeAuthorization(address addr) external {
        require(isValidVoter(msg.sender));
        authorizedMembers[addr] = false;
    }

    function isAuthorized(address addr) external view returns (bool) {
        if (existing(validManagers, addr)) {
            return true;
        }
        return authorizedMembers[addr];
    }

    function getManagers() external view returns (address[]) {
        return validManagers;
    }

    function removeValidMember(address addr) internal {
        uint index = validManagers.length;
        for (uint i = 0; i < validManagers.length; i++) {
            if (validManagers[i] == addr) {
                index = i;
                break;
            }
        }
        delete validManagers[index];
    }

    function addElement(address[] storage addrList, address addr) internal {
        for (uint i = 0; i < addrList.length; i++) {
            if (addrList[i] == address(0)) {
                addrList[i] = addr;
                return;
            }
        }
        addrList.push(addr);
    }

    function arraryLen(address[] storage addrList) internal view returns (uint) {
        uint length = 0;
        for (uint i = 0; i < addrList.length; i++) {
            if (addrList[i] != address(0)) {
                length++;
            }
        }
        return length;
    }

    function isValidVoter(address addr) internal view returns (bool) {
        return existing(validManagers, addr);
    }

    function existing(address[] storage addrList, address addr) internal view returns (bool) {
        for (uint i = 0; i < addrList.length; i++) {
            if (addrList[i] == addr) {
                return true;
            }
        }
        return false;
    }
}