// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Aegis} from "./Aegis.sol";
import {Account} from "../libraries/Account.sol";

abstract contract NexusContext is Aegis {
    address[] internal accountList;

    mapping(address => Account.Info) internal nexus;
    mapping(address => uint256) internal accountIndex;

    function _isMaster(address account) internal view override returns (bool) {
        return nexus[account].role == Account.Role.Master;
    }

    function _isMember(address account) internal view override returns (bool) {
        Account.Role role = nexus[account].role;
        return role == Account.Role.Master || role == Account.Role.Member;
    }
}