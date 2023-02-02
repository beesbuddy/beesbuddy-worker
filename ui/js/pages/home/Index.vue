<template>
    <app-title :title="subtitle" />
    <div class="p-5">
        <breadcrumb :links="links" />
        <card-title>
            <span class="flex-1 sm:mr-auto">
                Hello from {{ title }}
            </span>
        </card-title>
    </div>
</template>

<script>

import { ref } from 'vue';
import Layout from '../../common/Layout.vue';
import AppTitle from '../../common/AppTitle.vue';
import Breadcrumb from '../../common/Breadcrumb.vue';
import CardTitle from '../../common/CardTitle.vue';

export default {
    components: {
        AppTitle,
        Breadcrumb,
        CardTitle,
    },
    layout: Layout,
    props: {
        title: {
            type: String,
            default: () => 'appName',
        },
    },
    setup() {
        const subtitle = ref('Welcome');

        const links = ref([
            { name: subtitle }
        ]);

        return {
            subtitle,
            links
        };
    }
}
</script>
