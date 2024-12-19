<script setup lang="ts">
import { ref, onMounted } from 'vue';

const questions= ref<any[]>([]);
const isLoading = ref(true);

const props = withDefaults(defineProps<{ requestId: string; }>(), { requestId: '' })

onMounted(async () => {
  isLoading.value = true;
  const res = await fetch('http://localhost:8081/questions', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!res.ok) {
    console.error('Failed to fetch questions');
    return;
  }

  const questionsRaw = await res.json();
  questions.value = questionsRaw.map((q: any) => ({
    ...q,
    answer: q.type === 'text' ? '' : false,
  }));

  isLoading.value = false;
})

const submit = async () => {
  const res = await fetch('http://localhost:8081/answers', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      requestId: props.requestId,
      questions: questions.value,
    }),
  });

  if (!res.ok) {
    console.error('Failed to submit answers');
    return;
  }

  console.log('Answers submitted');
}
</script>

<template>
  <div>
    <p v-if="isLoading">Loading...</p>
    <template v-else>
      <div class="question" v-for="(q, idx) in questions">
        <template v-if="q.type === 'text'">
          <label :for="`question-${idx}`">{{ q.question }}</label>
          <br/>
          <input :id="`question-${idx}`" v-model="q.answer" />
        </template>

        <template v-else-if="q.type === 'boolean'">
          <p>{{ q.question }}</p>
          <label :for="`question-${idx}-true`">
            <input
                type="radio"
                :id="`question-${idx}-true`"
                :value="true"
                v-model="q.answer"
                :name="`question-${idx}`"
            />
            True
          </label>

          <label :for="`question-${idx}-false`">
            <input
                type="radio"
                :id="`question-${idx}-false`"
                :value="false"
                v-model="q.answer"
                :name="`question-${idx}`"
            />
            False
          </label>
        </template>
      </div>
      <button @click="submit">Submit</button>
    </template>
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
