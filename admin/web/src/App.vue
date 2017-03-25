<template>
    <v-app id="econfigApp" class="grey lighten-1" top-toolbar left-sidebar>
        <v-toolbar>
            <v-toolbar-side-icon class="hidden-md-and-up" @click.native.stop="sidebar4 = !sidebar4" />
            <v-toolbar-logo class="text-xs-right">Logo</v-toolbar-logo>
        </v-toolbar>
            <main>
                <v-sidebar v-model="sidebar4" height="auto">
                    <v-list dense>
                        <template v-for="(item,i) in itemGroup">
                            <v-list-group v-if="item.items">
                                <v-list-item slot="item">
                                    <v-list-tile :href="item.href" ripple>
                                        <v-list-tile-title v-text="item.title" />
                                            <v-list-tile-action>
                                                <v-icon>keyboard_arrow_down</v-icon>
                                            </v-list-tile-action>
                                     </v-list-tile>
                                 </v-list-item>
                                    <v-list-item v-for="(subItem,i) in item.items" :key="i">
                                        <v-list-tile :href="subItem.href" ripple>
                                            <v-list-tile-title v-text="subItem.title" />
                                        </v-list-tile>
                                    </v-list-item>
                               </v-list-group>
                               <v-subheader v-else-if="item.header" v-text="item.header" />
                               <v-divider v-else-if="item.divider" light />
                               <v-list-item v-else>
                                    <v-list-tile :href="item.href" ripple>
                                        <v-list-tile-title v-text="item.title" />
                                    </v-list-tile>
                                </v-list-item>
                            </template>
                        </v-list>
                    </v-sidebar>
                <v-content>
                <v-container fluid></v-container>
            </v-content>
        </main>
    </v-app>
</template>

<script>
    import config from './config'
    import project from './store/project'

    export default {
        name: 'app',
        data () {
            return {
                itemGroup: config.menu,
                projects: []
            }
        },
        created () {
            this.bindProjects()
        },

        methods: {
            bindProjects () {
                project.list(function (data) {
                    this.projects = data
                })
            }
        }
    }
</script>
<style>
    #app {}
</style>
