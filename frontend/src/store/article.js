import { defineStore } from "pinia";
import { computed, ref } from "vue";
// import { QueryAccounts, CreateAccounts, UpdateAccount, DeleteAccount } from '@/wailsjs/go/controller/DBController'


export const useArticleStore = defineStore("article", () => {
    const articleList = ref([]);
    

    return {
        articleList,
    }

});