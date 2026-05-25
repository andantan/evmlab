// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Account} from "../libraries/Account.sol";

interface INexus {
    error AlreadyRegisteredAccount(address account);
    error AlreadyApprovedAccount(address account);
    error NotRegisteredAccount(address account);
    error InvalidAccount(address account);
    error InvalidRole(Account.Role role);
    error CannotRemoveMaster();

    event AccountAdded(address indexed who, uint256 indexed when, Account.Role role, uint256 count);
    event AccountApproved(address indexed who, uint256 indexed when, Account.Role role, uint256 count);
    event AccountRemoved(address indexed who, uint256 indexed when, Account.Role role, uint256 count);
}
