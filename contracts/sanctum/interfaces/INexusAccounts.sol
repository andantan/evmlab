// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Account} from "../libraries/Account.sol";
import {INexus} from "./INexus.sol";

interface INexusAccounts is INexus {
    function register() external;
    function registerFor(address account) external;
    function deregister() external;
    function deregisterFor(address account) external;
    function approveRegister(address account) external;
    function getAccounts() external view returns (address[] memory);
    function accountCount() external view returns (uint256);
    function getAccountInfo(address account) external view returns (Account.Info memory);
}