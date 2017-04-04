<template>
    <div>
        <h2 style="margin: 10px">Configuration Items</h2>
        <Collapse style="margin-bottom: 20px;">
            <Panel name="1">
                Setter
                <div slot="content">
                    <Row>
                        <Col span="12">
                            <Form :model="selectedItem" :label-width="80">
                                <Form-item label="Key">
                                    <Select v-model="selectedKey"
                                        not-found-text="no data"
                                        placeholder="Select Item"
                                        clearable
                                        @on-change="handleSelectItem"
                                        filterable>
                                        <Option v-for="key in keys" :value="key" :key="key">{{ key }}</Option>
                                    </Select>
                                </Form-item>
                                <Form-item label="Value">
                                    <Input v-model="selectedItem.Value" placeholder="input value"></Input>
                                </Form-item>
                                <Form-item label="Tmpl">
                                    <Input v-model="selectedItem.TmplValue"></Input>
                                </Form-item>
                                <Form-item label="Store">
                                    <Input v-model="selectedItem.StoreValue"></Input>
                                </Form-item>
                                <Form-item>
                                    <Button type="primary" @click="handleSubmit('formValidate')">Save</Button>
                                    <Button type="ghost" @click="handleReset('formValidate')" style="margin-left: 8px">Reset</Button>
                                </Form-item>
                            </Form>
                        </Col>
                        <Col span="12">
                        </Col>
                    </Row>
                </div>
            </Panel>
        </Collapse>

        <h2 style="margin-bottom: 10px;">miss:<span style="color:red">{{ miss }}</span>/{{ total }}</h2>
        <Table border :context="self" :columns="table.columns" :data="table.data"></Table>
    </div>
</template>

<script>

import itemsStore from '../store/items'

export default {
    data () {
        return {
            self: this,
            selectedItem: {
                Key: '',
                Value: '',
                TmplValue: '',
                StoreValue: ''
            },
            table: {
                columns: [
                    {
                        title: 'Key',
                        key: 'Key',
                        width: 300
                    },
                    {
                        title: 'Value',
                        key: 'Value'
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
            storeItems: [],
            tmpItems: [],
            keys: [],
            selectedKey: ''
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
        bindItemValue: function (key, items) {
            if (key in items) {
                return items[key].Value
            }
            return ''
        }

    },

    methods: {
        handleSelectItem (item) {
            let my = this
            if (item) {
                my.selectedItem = {
                    Key: item,
                    Value: (item in my.items) ? my.items[item].Value : '',
                    TmplValue: my.tmpItems ? my.tmpItems[item].Value : '',
                    StoreValue: my.storeItems && (item in my.storeItems) ? my.storeItems[item].Value : ''
                }
            }
            return item
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
                    if (response.code === 0) {
                        my.tmpItems = response.data
                        let items = []
                        for (var prop in response.data) {
                            my.total ++
                            my.keys.push(prop)
                            let item = {}
                            item['Key'] = prop
                            if (prop in my.items) {
                                item['Value'] = my.items[prop].Value
                            } else {
                                item['Value'] = '-'
                                item['cellClassName'] = {
                                    'Value': 'miss'
                                }
                            }
                            item['Comment'] = response.data[prop].Comment
                            items.push(item)
                        }
                        my.table.data = items
                    } else {
                        my.$Message.error(response.msg, 10)
                    }
                })
            })
            itemsStore.listOfStore(my.appName, function (response) {
                if (response.code > 0) {
                    my.$Message.error(response.msg, 10)
                    return
                }
                my.storeItems = response.data
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
