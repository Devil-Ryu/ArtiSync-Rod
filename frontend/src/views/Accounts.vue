<template>
    <div class="m-3 h-full ">
        <Toast />
        <Card class="h-full">
            <template #content>
                <Toolbar class="mb-4">
                    <template #start>
                        <Button label="新增账户" icon="pi pi-plus" severity="success" class="mr-2" @click="newAccount" />
                        <!-- <Button label="测试授权" icon="pi pi-upload" severity="help" @click="()=>{console.log('refresh platforms'); RefreshPlatforms().then(res=>console.log('res:', res))}" /> -->
                    </template>

                    <template #end>
                        <FileUpload mode="basic" accept="image/*" :maxFileSize="1000000" label="Import" chooseLabel="导入"
                            class="mr-2 inline-block" :disabled="true"/>
                        <Button label="导出" icon="pi pi-upload" severity="help" @click="exportCSV($event)" :disabled="true"/>
                    </template>
                </Toolbar>
                <DataTable paginator :rows="10" :value="accountsStore.accounts" v-model:filters="dataFilters" filterDisplay="menu"
                    editMode="cell" size="small">

                    <!-- <Column field="PlatformKey" header="平台Key"></Column> -->
                    <Column field="PlatformAlias" header="平台名称"></Column>
                    <Column field="Username" header="用户名"></Column>
                    <Column field="LoginType" header="登录方式">
                        <template #body="{ data, field }">
                            {{  accountsStore.loginTypesMap[data[field]]  }}
                        </template>
                    </Column>
                    <Column field="Disabled" header="状态">
                        <template #body="{ data, field }">
                            <Tag :severity="data[field] ? 'danger' : 'success'">{{ data[field] ? '禁用' : '启用' }}</Tag>
                        </template>
                    </Column>
                    <Column :exportable="false" class="w-6rem">
                        <template #body="slotProps">
                            <Button icon="pi pi-pencil" text rounded @click="editAccount(slotProps.data)" />
                            <Button icon="pi pi-trash" text rounded severity="danger"
                                @click="deleteAccount(slotProps.data)" />
                        </template>
                    </Column>
                </DataTable>
            </template>
        </Card>
    </div>

    <Dialog v-model:visible="accountDialog" :style="{ width: '450px' }" header="账户详情"   :modal="true" class="p-fluid">
        <template #header>
            <div class="flex justify-content-center w-full col-offset-1">
                <div class="flex align-items-center justify-content-center font-bold text-xl">{{ account.Username ? '编辑账户' : '新增账户' }}</div>

            </div>
        </template>
        
        <TabView 
            v-model:activeIndex="activeTabIndex"
        >
            <TabPanel header="跳转授权">
                <div class="field">
                    <label for="PlatformKey">平台名称</label>
                    <Dropdown id="PlatformKey" v-model="account.PlatformKey" :options="platformStore.platforms"
                        optionLabel="Alias" optionValue="Key" placeholder="请选择平台">
                    </Dropdown>
                </div>
                <div class="field">
                    <Button label="点击授权" icon="pi pi-external-link" class="mr-2" :disabled="!account.PlatformKey" @click="authAccount" />
                </div>
            </TabPanel>
            <TabPanel header="账户编辑">
                <div class="field">
                    <label for="PlatformKey">平台名称</label>
                    <Dropdown id="PlatformKey" v-model="account.PlatformKey" :options="platformStore.platforms"
                        optionLabel="Alias" optionValue="Key" placeholder="请选择">
                    </Dropdown>
                </div>
                <div class="field">
                    <label for="LoginType">登录方式</label>
                    <Dropdown id="LoginType" v-model="account.LoginType" :options="accountsStore.loginTypes" optionLabel="label" optionValue="value"
                        placeholder="请选择">
                    </Dropdown>
                </div>
                <div class="field">
                    <label for="Username">用户名</label>
                    <InputText id="Username" v-model.trim="account.Username" required="true" autofocus />
                </div>
                <div class="field" v-if="account.LoginType !== 'AUTH_REDIRECT'">
                    <label for="Password">密码</label>
                    <InputText id="Password" v-model="account.Password" required="true" />
                </div>
            </TabPanel>
        </TabView>

        <template #footer>
                    <Button label="取消" icon="pi pi-times" text />
                    <Button label="保存" icon="pi pi-check" text @click="saveAccount" />
                </template>
    </Dialog>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { QueryAccounts, CreateAccounts, UpdateAccount, DeleteAccount } from '@/wailsjs/go/controller/DBController'
import { RefreshPlatforms } from '@/wailsjs/go/application/ATApp'
import { usePlatformsStore } from '@/src/store/platform.js'
import { useToast } from 'primevue/usetoast';
import { useAccountsStore } from '@/src/store/accounts.js'
const accountsStore = useAccountsStore()


const toast = useToast();
const platformStore = usePlatformsStore()
const accountDialog = ref(false);
const account = ref({});
const activeTabIndex = ref(0);

const dataFilters = ref({
})


const newAccount = () => {
    activeTabIndex.value = 0;
    account.value = {};
    accountDialog.value = true;
};

const editAccount = (act) => {
    activeTabIndex.value = 1;
    account.value = { ...act };
    accountDialog.value = true;
    console.log("edit account", act);
};

const saveAccount = () => {
    if (account.value.ID) {
        // update account
        UpdateAccount(account.value).then(res => {
            console.log("update account", res);
            let index = accountsStore.accounts.findIndex(item => item.ID === account.value.ID);
            accountsStore.accounts[index] = res;
            toast.add({ severity: 'success', summary: '更新成功', detail: res.Username + "已经更新", life: 3000 });
            accountDialog.value = false;
            account.value = {};
        }).catch(err => {
            console.log("update account error", err);
            toast.add({ severity: 'error', summary: '更新失败', detail: err, life: 3000 });
        })
    } else {
        // create account
        let platformAlias = platformStore.platforms.find(item => item.Key === account.value.PlatformKey).Alias;
        account.value.PlatformAlias = platformAlias;  // 添加平台名称
        console.log("create account", account.value);
        CreateAccounts([account.value]).then(res => {
            console.log("create account", res[0]);
            accountsStore.accounts.push(res[0]);
            toast.add({ severity: 'success', summary: '创建成功', detail: res[0].Username + " 账户创建成功", life: 3000 });
            accountDialog.value = false;
            account.value = {};
        }).catch(err => {
            console.log("create account error", err);
            toast.add({ severity: 'error', summary: '创建失败', detail: err, life: 3000 });
        })
    }
};

const deleteAccount = (act) => {
    console.log("delete account", act);
    DeleteAccount(act).then(res => {
        console.log("delete account", res);
        let index = accountsStore.accounts.findIndex(item => item.ID === act.ID);
        accountsStore.accounts.splice(index, 1);
        toast.add({ severity: 'success', summary: '删除成功', detail: act.Username + " 账户删除成功", life: 3000 });
    }).catch(err => {
        console.log("delete account error", err);
        toast.add({ severity: 'error', summary: '删除失败', detail: err, life: 3000 });
    })
};

const authAccount = () => {
    console.log("auth account", account.value);
    let platform = platformStore.platforms.find(item => item.Key === account.value.PlatformKey);
    platform.login().then(res => {
        if (res === null) {
            toast.add({ severity: 'warn', summary: '授权取消', detail: "返回账户为: null", life: 3000 });
            return;
        }
        account.value = {
            ID: account.value.ID,
            PlatformKey: res.PlatformKey,
            PlatformAlias: res.PlatformAlias,
            Username: res.Username,
            LoginType: res.LoginType,
            Password: res.Password,
            Cookies: res.Cookies,
            Disabled: false,
        }
        toast.add({ severity: 'success', summary: '授权成功', detail: res.Username + " 账户授权成功", life: 3000 });
        // 没有跳转TODO
        activeTabIndex.value = 1;
    }).catch(err => {
        console.log("auth account error", err);
        toast.add({ severity: 'error', summary: '授权失败', detail: err, life: 3000 });
    })
};


// onMounted(() => {
//     console.log('mounted');
//     QueryAccounts().then(res => {
//         console.log("res", res);
//         accountsStore.accounts = res;
//     })
// });

</script>