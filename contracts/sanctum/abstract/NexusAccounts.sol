// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SanctumContext} from "./SanctumContext.sol";
import {INexusAccounts} from "../interfaces/INexusAccounts.sol";
import {Account} from "../libraries/Account.sol";

abstract contract NexusAccounts is SanctumContext, INexusAccounts {
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
        if (!isAccount[account]) {
            revert NotRegisteredAccount(account);
        }

        return Account.Info({
            account: account,
            role: roles[account],
            registeredBlock: registeredBlock[account]
        });
    }

    function _addAccount(address account, Account.Role role) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (isAccount[account]) {
            revert AlreadyRegisteredAccount(account);
        }

        if (!Account.isValidRole(role)) {
            revert InvalidRole(role);
        }

        isAccount[account] = true;
        roles[account] = role;
        accountIndex[account] = accountList.length;
        registeredBlock[account] = block.number;
        accountList.push(account);

        emit AccountAdded(account, block.number, role, accountList.length);
    }

    function _approveAccount(address account, Account.Role role) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (!isAccount[account]) {
            revert NotRegisteredAccount(account);
        }

        if (!Account.isValidRole(role)) {
            revert InvalidRole(role);
        }

        if (!Account.isPending(roles[account])) {
            revert AlreadyApprovedAccount(account);
        }

        roles[account] = role;

        emit AccountApproved(account, block.number, role, accountList.length);
    }

    function _removeAccount(address account) internal {
        if (account == address(0)) {
            revert InvalidAccount(account);
        }

        if (!isAccount[account]) {
            revert NotRegisteredAccount(account);
        }

        if (Account.isMaster(roles[account])) {
            revert CannotRemoveMaster();
        }

        Account.Role role = roles[account];
        uint256 index = accountIndex[account];
        uint256 lastIndex = accountList.length - 1;

        if (index != lastIndex) {
            address lastAccount = accountList[lastIndex];
            accountList[index] = lastAccount;
            accountIndex[lastAccount] = index;
        }

        accountList.pop();

        delete accountIndex[account];
        delete isAccount[account];
        delete roles[account];
        delete registeredBlock[account];

        emit AccountRemoved(account, block.number, role, accountList.length);
    }
}