<script setup lang="ts">
import {computed, ref} from 'vue';
import ModuleRenderer from "./components/ModuleRenderer.vue";

const eventSource = ref<EventSource | null>(null);
const userId = ref<string>('');
const requestId = ref<string>('');
const eventsLog = ref<any[]>([]);

const reverseEventsLog = computed(() => eventsLog.value.reverse());

const newRequest = async () => {
  const response = await fetch('http://localhost:8080/start', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ userId: userId.value }),
  });
  if (!response.ok) {
    eventsLog.value.push('Failed to create request' + response.statusText);
    return;
  }
  const data = await response.json();
  requestId.value = data.requestId;
  eventsLog.value.push(`Request created: ${requestId.value}`);

  resetEventSource();
};

const resetEventSource = () => {
  if (eventSource.value) {
    eventSource.value.close();
  }
  initEventSource();
}

const module = ref({
  url: '',
  tag: '',
});
const initEventSource = () => {
  const orchestratorApiUrl = `http://localhost:8080/events?requestId=${requestId.value}`;

  eventSource.value = new EventSource(orchestratorApiUrl);
  eventSource.value.onmessage = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    if (data.module) {
      module.value = { url: data.module.moduleURL, tag: data.module.moduleTag };
    }
    eventsLog.value = [...eventsLog.value, data];
  };
}
</script>

<template>
  <div class="container">
    <div class="container-left">
      <header class="container-header">
        <label for="userID">User ID:</label>
        <input v-model="userId" id="userID">
        <button @click="newRequest">Request</button>
      </header>
      <main class="module-target">
        <ModuleRenderer :module-url="module.url" :module-tag="module.tag" :request-id="requestId"/>
      </main>
    </div>
    <aside class="logger">
      <ul class="log-list">
        <li v-for="log in reverseEventsLog" :key="log">
          <p>{{ log }}</p>
        </li>
      </ul>
    </aside>
  </div>
</template>

<style scoped lang="scss">
.container {
  display: flex;
  gap: 1rem;
  width: 100vw;
  height: 100vh;
  .container-left {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 1rem;
    .container-header { display: flex; gap: 0.5rem; }
  }
  .module-target { flex-grow: 1; }
  .logger {
    width: 30vw;
    border: 1px solid #ccc;
    border-radius: 2px;
    padding: 5px;
    overflow-y: auto;
    .log-list {
      list-style: none;
      display: flex;
      flex-direction: column;
      gap: 0.2rem;
      padding: 0;
      margin: 0;
      li {
        font-family: "Droid Sans Mono";
        border-bottom: 1px solid #999;
        overflow: hidden;
        white-space: break-spaces;
        word-break: break-all;
      }
    }
  }
}
</style>
