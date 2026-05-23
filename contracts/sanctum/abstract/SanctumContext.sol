// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Account} from "../libraries/Account.sol";

abstract contract SanctumContext {
    address public master;
    address[] internal accountList;

    mapping(address => bool) internal isAccount;
    mapping(address => Account.Role) internal roles;
    mapping(address => uint256) internal accountIndex;
    mapping(address => uint256) internal registeredBlock;

    error UnauthorizedContext(address caller);

    modifier onlyMaster() {
        if (msg.sender != master) {
            revert UnauthorizedContext(msg.sender);
        }

        _;
    }

    modifier onlyMember() {
        if (!isAccount[msg.sender] || !Account.isMember(roles[msg.sender])) {
            revert UnauthorizedContext(msg.sender);
        }

        _;
    }
}