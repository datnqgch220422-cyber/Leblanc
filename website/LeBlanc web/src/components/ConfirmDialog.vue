<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Confirm action",
  },
  message: {
    type: String,
    default: "Are you sure?",
  },
  confirmText: {
    type: String,
    default: "Confirm",
  },
  cancelText: {
    type: String,
    default: "Cancel",
  },
  pending: {
    type: Boolean,
    default: false,
  },
  danger: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:modelValue", "confirm", "cancel"]);

const close = () => {
  if (props.pending) return;
  emit("update:modelValue", false);
  emit("cancel");
};

const confirm = () => {
  if (props.pending) return;
  emit("confirm");
};
</script>

<template>
  <div v-if="modelValue" class="confirm-overlay" @click.self="close">
    <div class="confirm-card" role="dialog" aria-modal="true">
      <h3>{{ title }}</h3>
      <p>{{ message }}</p>
      <div class="confirm-actions">
        <button class="confirm-btn" @click="close" :disabled="pending">
          {{ cancelText }}
        </button>
        <button
          class="confirm-btn"
          :class="{ 'confirm-btn--danger': danger }"
          @click="confirm"
          :disabled="pending"
        >
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: grid;
  place-items: center;
  z-index: 1000;
  padding: 16px;
}

.confirm-card {
  width: min(92vw, 440px);
  background: #ffffff;
  border-radius: 12px;
  padding: 18px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
}

.confirm-card h3 {
  margin: 0 0 6px;
}

.confirm-card p {
  margin: 0 0 14px;
  color: rgba(0, 0, 0, 0.75);
}

.confirm-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.confirm-btn {
  border: 1px solid rgba(0, 0, 0, 0.15);
  background: #ffffff;
  border-radius: 8px;
  padding: 8px 12px;
  font-weight: 700;
  cursor: pointer;
}

.confirm-btn--danger {
  border-color: #9d3412;
  color: #9d3412;
}

.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
