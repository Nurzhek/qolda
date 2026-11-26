<template>
  <div>
    <ul class="grid grid-cols-3 gap-2 sm:grid-cols-6">
      <li
        v-for="(url, index) in photos"
        :key="url"
        class="relative aspect-square overflow-hidden rounded-[3px] border border-rule bg-paper-deep"
      >
        <img
          :src="url"
          alt=""
          class="size-full object-cover"
        >
        <button
          type="button"
          class="absolute top-1 right-1 grid size-7 place-items-center rounded-[2px] border border-rule bg-sheet text-ink hover:bg-paper-deep"
          :aria-label="t('listing.form.removePhoto', { n: index + 1 })"
          @click="remove(index)"
        >
          <PhX
            :size="14"
            weight="bold"
            aria-hidden="true"
          />
        </button>
      </li>

      <li
        v-for="i in uploading"
        :key="`up-${i}`"
        class="aspect-square"
      >
        <USkeleton class="size-full rounded-[3px]" />
        <span class="sr-only">{{ t('listing.form.uploading') }}</span>
      </li>

      <li v-if="canAdd" class="aspect-square">
        <label
          class="flex size-full cursor-pointer flex-col items-center justify-center gap-1 rounded-[3px] border border-dashed border-rule-strong text-meta text-graphite hover:border-ink hover:text-ink has-[input:focus-visible]:outline-2 has-[input:focus-visible]:outline-moss-500"
        >
          <PhPlus :size="20" aria-hidden="true" />
          {{ t('listing.form.addPhoto') }}
          <input
            type="file"
            :accept="ACCEPT.join(',')"
            multiple
            class="sr-only"
            @change="onPick"
          >
        </label>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { Executor } from 'shared/api/executor';
import { useLoading } from 'shared/lib/composables/useLoading';
import { useI18n } from 'vue-i18n';
import { FileRepository } from '../api/file.repository';

const props = defineProps<{
  max: number;
}>();

const photos = defineModel<string[]>({ default: () => [] });

const ACCEPT = ['image/jpeg', 'image/png', 'image/webp'];
const MAX_BYTES = 8 * 1024 * 1024;

const { t } = useI18n();
const toast = useToast();
const loadingState = useLoading();
const uploading = ref(0);

const canAdd = computed(() => photos.value.length + uploading.value < props.max);

function remove(index: number) {
  photos.value = photos.value.filter((_, i) => i !== index);
}

function reject(code: 'file_too_large' | 'file_type_not_allowed') {
  toast.add({ title: t(`errors.${code}`), color: 'error' });
}

async function upload(file: File) {
  uploading.value += 1;
  const result = await Executor.run({
    request: FileRepository.createFile(file),
    loadingState,
  });
  uploading.value -= 1;
  if (result.ok)
    photos.value = [...photos.value, result.data.url];
}

async function onPick(event: Event) {
  const input = event.target as HTMLInputElement;
  const free = props.max - photos.value.length - uploading.value;
  const files = Array.from(input.files ?? []).slice(0, Math.max(free, 0));
  input.value = '';

  const valid = files.filter((file) => {
    if (!ACCEPT.includes(file.type)) {
      reject('file_type_not_allowed');
      return false;
    }
    if (file.size > MAX_BYTES) {
      reject('file_too_large');
      return false;
    }
    return true;
  });

  // Photo order = selection order: upload one at a time
  for (const file of valid)
    await upload(file);
}
</script>
