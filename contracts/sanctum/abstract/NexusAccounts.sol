// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {NexusContext} from "./NexusContext.sol";
import {INexusAccounts} from "../interfaces/INexusAccounts.sol";
import {Account} from "../libraries/Account.sol";

abstract contract NexusAccounts is NexusContext, INexusAccounts {
    function register() external {
        _addAccount(msg.sender, Account.Role.Pending);
    }

    function registerFor(address account) external onlyMaster {
        _addAccount(account, Account.Role.Member);
    }

    function approveRegister(address account) external onlyMaster {
        _approveAccount(account, Account.Role.Member);
    }

    function deregister() external {
        _removeAccount(msg.sender);
    }

    function deregisterFor(address account) external onlyMaster {
        _removeAccount(account);
    }

    function getAccounts() external view returns (address[] memory) {
        return accountList;
    }

    function accountCount() external view returns (uint256) {
        return accountList.length;
    }

    function getAccountInfo(address account) external view returns (Account.Info memory) {
        if (nexus[account].addr == address(0)) {
            revert NotRegisteredAccount(account);
        }

        return nexus[account];
    }

    function _addAccount(address account, Account.Role role) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (role == Account.Role.None) {
            revert InvalidRole(role);
        }

        if (nexus[account].addr != address(0)) {
            revert AlreadyRegisteredAccount(account);
        }

        if (uint8(role) > uint8(type(Account.Role).max)) {
            revert InvalidRole(role);
        }

        nexus[account] = Account.Info({
            addr: account,
            role: role,
            registeredBlock: block.number
        });
        accountIndex[account] = accountList.length;
        accountList.push(account);

        emit AccountAdded(account, block.number, role, accountList.length);
    }

    function _approveAccount(address account, Account.Role role) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (nexus[account].addr == address(0)) {
            revert NotRegisteredAccount(account);
        }

        if (uint8(role) > uint8(type(Account.Role).max)) {
            revert InvalidRole(role);
        }

        if (nexus[account].role != Account.Role.Pending) {
            revert AlreadyApprovedAccount(account);
        }

        nexus[account].role = role;

        emit AccountApproved(account, block.number, role, accountList.length);
    }

    function _removeAccount(address account) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (nexus[account].addr == address(0)) {
            revert NotRegisteredAccount(account);
        }

        if (nexus[account].role == Account.Role.Master) {
            revert CannotRemoveMaster();
        }

        Account.Role role = nexus[account].role;
        uint256 index = accountIndex[account];
        uint256 lastIndex = accountList.length - 1;

        if (index != lastIndex) {
            address lastAccount = accountList[lastIndex];
            accountList[index] = lastAccount;
            accountIndex[lastAccount] = index;
        }

        accountList.pop();

        delete accountIndex[account];
        delete nexus[account];

        emit AccountRemoved(account, block.number, role, accountList.length);
    }
}