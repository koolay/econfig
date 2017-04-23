<template>
    <div>
        <h2 style="margin: 10px">Configuration Items</h2>
        <Collapse style="margin-bottom: 20px;">
            <Panel name="1">
                Setter
                <div slot="content">
                    <Row>
                        <Col span="12">
                            <Form :model="formItem" ref="formItem" :label-width="80">
                                <Form-item label="Key">
                                    <Select v-model="formItem.Key"
                                        not-found-text="no data"
                                        placeholder="Select Item"
                                        clearable
                                        @on-change="handleSelectItem"
                                        filterable>
                                        <Option v-for="key in keys" :value="key" :label="key" :key="key">
                                            <p>{{ key }}</p>
                                            <p style="color:#ccc">{{ tmpItems[key].Comment}}</p>
                                        </Option>
                                    </Select>
                                </Form-item>
                                <Form-item label="Dest Value">
                                    <Input type="textarea" :autosize="true" v-model="formItem.Value" placeholder="input value"></Input>
                                </Form-item>
                                <Form-item label="Example">
                                    <div>{{ selectedItem.TmplValue }}</div>
                                </Form-item>
                                <Form-item label="Store">
                                    <div>{{ selectedItem.StoreValue }}</div>
                                </Form-item>
                                <Form-item>
                                    <Button type="primary" @click="handleSave('formItem')">Save</Button>
                                    <Button type="ghost" @click="handleReset('formItem')" style="margin-left: 8px">Reset</Button>
                                    <Button type="primary" @click="handleSync()" style="margin-left: 8px">Sync Dest file</Button>
                                </Form-item>
                            </Form>
                        </Col>
                        <Col span="1">
                        &nbsp;
                        </Col>
                        <Col span="11">
                        <h3 style="color:#00cc66" v-if="syncResult.Logs.length>0 && syncResult.Success">Success</h3>
                        <h3 style="color:#ff3300" v-if="syncResult.Logs.length>0 && !syncResult.Success">Failure</h3>
                            <ul class="sync-log">
                                <li v-for="(log, index) in syncResult.Logs">{{ log }}</li>
                            </ul>
                        </Col>
                    </Row>
                </div>
            </Panel>
        </Collapse>

        <h2 style="margin-bottom: 10px;">miss:<span style="color:red">{{ miss }}</span>/{{ total }}</h2>
        <Table border :context="self" :columns="table.columns" :data="table.data" @on-row-click="handleRowClick"></Table>
    </div>
</template>

<script>

import itemsStore from '../store/items'
import appStore from '../store/app'

export default {
    data () {
        return {
            self: this,
            syncResult: {
                Logs: [],
                Success: true
            },
            selectedItem: {
                Key: '',
                Value: '',
                TmplValue: '',
                StoreValue: '',
                Comment: ''
            },
            formItem: {
                Key: '',
                Value: ''
            },
            table: {
                columns: [
                    {
                        title: 'Key',
                        key: 'Key',
                        width: 300
                    },
                    {
                        title: 'Dest Value',
                        key: 'Value'
                    },
                    {
                        title: 'Store Value',
                        key: 'StoreValue'
                    },
                    {
                        title: 'Comment',
                        key: 'Comment',
                        width: 250
                    }
                ],
                data: []
            },
            appName: this.$route.params.app,
            items: [],
            storeItemsMap: [],
            tmpItems: [],
            keys: []
        }
    },

    computed: {
        miss: function () {
            return this.total - Object.keys(this.items).length
        },
        total: function () {
            return Object.keys(this.tmpItems).length
        }
    },
    created () {
        this.bindItems()
    },

    filters: {

        shortSelectComment: function (comment) {
            if (comment && comment.length > 50) {
                return comment + '...'
            }
            return comment
        }
    },

    methods: {
        handleSelectItem (item) {
            let my = this
            if (item) {
                let storeVal = ''
                if (my.storeItemsMap && Object.prototype.hasOwnProperty.call(my.storeItemsMap, item)) {
                    storeVal = my.storeItemsMap[item]
                }
                my.selectedItem = {
                    Key: item,
                    Value: (item in my.items) ? my.items[item].Value : '',
                    TmplValue: my.tmpItems ? my.tmpItems[item].Value : '',
                    StoreValue: storeVal,
                    Comment: my.tmpItems ? my.tmpItems[item].Comment : ''
                }
                my.formItem = {...my.selectedItem}
            }
            return item
        },
        handleRowClick (itemObj) {
            this.formItem.Key = itemObj.Key
        },
        handleSave (name) {
            let my = this
            let key = this.formItem.Key
            let value = this.formItem.Value
            itemsStore.saveItem(my.appName, key, value, response => {
                if (response.code > 0) {
                    my.$Message.error(response.msg, 10)
                    return
                }
                my.storeItemsMap[key] = value
                my.selectedItem.StoreValue = value
                my.$Message.success('Successfully saved')
            })
        },
        handleReset (name) {
            this.formItem = {...this.selectedItem}
        },
        handleSync () {
            let my = this
            appStore.execSync(this.appName, response => {
                if (response.code > 0) {
                    my.$Message.error(response.msg, 10)
                    return
                }
                my.syncResult = response.data
            })
        },
        bindItems () {
            let my = this
            itemsStore.list(my.appName, function (response) {
                my.items = response.data
                if (response.code > 0) {
                    my.$Message.error(response.msg, 10)
                    return
                }

                itemsStore.listOfTmp(my.appName, function (response) {
                    if (response.code > 0) {
                        my.$Message.error(response.msg, 10)
                    }
                    my.tmpItems = response.data

                    itemsStore.listOfStore(my.appName, function (storeResponse) {
                        if (storeResponse.code > 0) {
                            my.$Message.error(storeResponse.msg, 10)
                            return
                        }
                        my.storeItemsMap = storeResponse.data
                        let items = []
                        for (var prop in response.data) {
                            my.total ++
                            my.keys.push(prop)
                            let item = {}
                            let cellsCss = {}
                            item['Key'] = prop
                            if (prop in my.items) {
                                item['Value'] = my.items[prop].Value
                            } else {
                                item['Value'] = '-'
                                cellsCss['Value'] = 'miss'
                            }
                            if (Object.prototype.hasOwnProperty.call(storeResponse.data, prop)) {
                                item['StoreValue'] = storeResponse.data[prop]
                            } else {
                                item['StoreValue'] = '-'
                                cellsCss['StoreValue'] = 'miss'
                            }
                            item['cellClassName'] = cellsCss
                            item['Comment'] = response.data[prop].Comment
                            items.push(item)
                        }
                        my.table.data = items
                    })
                })
            })
        }
    }
}
</script>
<style>
    .ivu-table .miss {
        background-color: #f60;
        color: #fff;
    }
</style>
