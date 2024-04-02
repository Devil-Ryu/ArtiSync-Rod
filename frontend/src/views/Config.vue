
<template>
        <TreeTable :value="nodes" >
            <Column field="name" header="Name" expander></Column>
            <Column field="size" header="Size"></Column>
            <Column field="data" header="Data"></Column>
        </TreeTable>
</template>

<script setup>
import { ref } from 'vue';



const nodes = ref([
{
    data: {name:'ddd', size: '123', data:"123456"},
    
    children: [
        {
           
            data: {name:'ddd', size: '123', data:"123456"},
            children: [
                { key: '0-0-0', label: 'Expenses.doc', icon: 'pi pi-fw pi-file', data: {name:'ddd', size: '123', data:"123456"} },
                { key: '0-0-1', label: 'Resume.doc', icon: 'pi pi-fw pi-file', data: {name:'ddd', size: '123', data:"123456"} }
            ]
        },
       
    ]
},
]);
</script>
