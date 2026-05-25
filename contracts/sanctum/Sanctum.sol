// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {NexusAccounts} from "./abstract/NexusAccounts.sol";
import {TreasuryNative} from "./abstract/TreasuryNative.sol";
import {Account} from "./libraries/Account.sol";

contract Sanctum is NexusAccounts, TreasuryNative {
    constructor() {
        master = msg.sender;
        _addAccount(msg.sender, Account.Role.Master);
    }
}