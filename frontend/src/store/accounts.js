import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { QueryAccounts, CreateAccounts, UpdateAccount, DeleteAccount } from '@/wailsjs/go/controller/DBController'


export const useAccountsStore = defineStore("accounts", () => {
    const accounts = ref([]);

    const loginTypes = ref([
        { label: "手动录入", value: "AUTH_MANUAL" },
        { label: "批量导入", value: "AUTH_BATCH" },
        { label: "跳转授权", value: "AUTH_REDIRECT" },
    ]);

    const loginTypesMap = computed(() => {
        const map = {};
        loginTypes.value.forEach(loginType => {
            map[loginType.value] = loginType.label;

        });
        return map;
    });

    const platformsCount = computed(() => {
        const platforms = {};
        accounts.value.forEach(account => {
            if (platforms[account.PlatformKey]) {
                platforms[account.PlatformKey] += 1;
            } else {
                platforms[account.PlatformKey] = 1;
            }
        });
        return platforms;
    });

    function RefreshAccounts() {
        console.log("Refreshing accounts");
        QueryAccounts().then(res => {
            console.log(res);
            accounts.value = res;
        })
    }

    return {
        accounts,
        loginTypes,
        loginTypesMap,
        platformsCount,
        RefreshAccounts
    }

});