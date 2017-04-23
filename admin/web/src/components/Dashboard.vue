<template>
    <div>
        <h2 style="margin: 10px">App List</h2>
        <Card style="width:45%; float: left; margin: 10px" v-for="(app, index) in apps" :key="app.Name">
            <p slot="title">
                <Icon type="ios-film-outline"></Icon>
                {{ app.Name}}
            </p>
            <!-- <a href="#" style="margin-left: 10px;" slot="extra" @click.prevent="changeLimit">
                <Tooltip content="modify" placement="top">
                    <Icon type="settings"></Icon>
                </Tooltip>
            </a> -->

            <router-link :to="{ name: 'items', params: { app: app.Name } }" style="margin-left: 10px;" slot="extra" @click.prevent="changeLimit">
                <Tooltip content="settings" placement="top">
                    <Icon type="navigate"></Icon>
                </Tooltip>
            </router-link>
            <row>
                <Col span="4">Root</Col>
                <Col span="20">{{ app.Root }}</Col>
            </row>
            <row>
                <Col span="4">Dest</Col>
                <Col span="20">{{ app.Dest}}</Col>
            </row>
            <row>
                <Col span="4">Tmpl</Col>
                <Col span="20">{{ app.Tmpl }}</Col>
            </row>
            <row>
                <Col span="4">Prefix</Col>
                <Col span="20">{{ app.Prefix }}</Col>
            </row>
            <row>
                <Col span="4">Cmd</Col>
                <Col span="20">{{ app.Cmd }}</Col>
            </row>
        </Card>
    </div>
</template>

<script>

    import config from '../config'
    import app from '../store/app'

    export default {
        data () {
            return {
                itemGroup: config.menu,
                apps: []
            }
        },
        created () {
            this.bindProjects()
        },

        methods: {
            bindProjects () {
                let my = this
                app.list(function (response) {
                    my.apps = response.data
                })
            }
        }
    }
</script>
