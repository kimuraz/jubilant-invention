<script setup lang="ts">
import {ref, watch} from "vue";

interface Props {
  moduleUrl: string;
  moduleTag: string;
  requestId: string;
}

const props = defineProps<Props>();
const moduleContainer = ref<HTMLElement | null>(null);

watch(props, async (newProps) => {
  if (newProps.moduleUrl && newProps.moduleTag && newProps.requestId) {
    await import(newProps.moduleUrl);
    const moduleElement = document.createElement(newProps.moduleTag);
    moduleElement.setAttribute('request-id', newProps.requestId);
    moduleContainer.value.appendChild(moduleElement);
  }
}, { immediate: true, deep: true });
</script>

<template>
  <div ref="moduleContainer"/>
</template>

<style scoped>

</style>