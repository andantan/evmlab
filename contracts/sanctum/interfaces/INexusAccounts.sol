// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Account} from "../libraries/Account.sol";

interface INexusAccounts {
    error AlreadyRegisteredAccount(address account);
    error AlreadyApprovedAccount(address account);
    error NotRegisteredAccount(address account);
    error InvalidAccount(address account);
    error InvalidRole(Account.Role role);
    error CannotRemoveMaster();

    event AccountAdded(
        address indexed who,
        uint256 indexed when,
        Account.Role indexed role,
        uint256 count
    );

    event AccountApproved(
        address indexed who,
        uint256 indexed when,
        Account.Role indexed role,
        uint256 count
    );

    event AccountRemoved(
        address indexed who,
        uint256 indexed when,
        Account.Role indexed role,
        uint256 count
    );

    // caller registers themselves as Pending — requires master approval to become Member
    function register() external;
    // master registers an account directly as Member
    function registerFor(address account) external;
    // caller deregisters themselves
    function deregister() external;
    // master deregisters any account
    function deregisterFor(address account) external;
    function approveRegister(address account) external;
    function getAccounts() external view returns (address[] memory);
    function accountCount() external view returns (uint256);
    function getAccountInfo(address account) external view returns (Account.Info memory);
}