<script setup lang="ts">
import { ref } from "vue";

interface Props {
  custom?: boolean;
  useDefaultBackground?: boolean;
  clickOutsideToClose?: boolean;
}
const props = defineProps<Props>();
const open = ref(false);

let mousedownLocation: { x: number; y: number } | null = null;

function recordMouseLocation(e: MouseEvent) {
  if (!props.clickOutsideToClose) {
    return;
  }
  mousedownLocation = {
    x: e.clientX,
    y: e.clientY,
  };
}

function closeModalIfNotDragging(e: MouseEvent) {
  if (!props.clickOutsideToClose || !mousedownLocation) {
    return;
  }
  const xmove = e.clientX - mousedownLocation.x;
  const ymove = e.clientY - mousedownLocation.y;
  mousedownLocation = null;
  if (xmove * xmove + ymove * ymove > 25) {
    return;
  }
  closeModal();
}

function openModal() {
  open.value = true;
}

function closeModal() {
  open.value = false;
}
</script>

<template>
  <slot name="trigger" :open="openModal">
    <button @click="openModal">Open Modal</button>
  </slot>
  <slot v-if="custom && !useDefaultBackground" :close="closeModal"> </slot>
  <!-- default background -->
  <div
    v-else-if="open"
    fixed
    inset-0
    z-50
    flex
    items-center
    justify-center
    bg-black
    bg-opacity-50
    @mousedown="recordMouseLocation"
    @mouseup="closeModalIfNotDragging"
  >
    <div flex items-center justify-center @mousedown.stop @mouseup.stop @click.stop>
      <slot :close="closeModal">
        <div relative bg-white p-4 rounded-lg max-w-160>
          <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="closeModal" />

          <h3 text-center text-base font-semibold text-gray-900>Lorem ipsum</h3>

          <p text-sm text-gray-500>
            Lorem ipsum dolor sit, amet consectetur adipisicing elit. Ea vero voluptatibus laborum
            laudantium accusamus culpa quas qui ut quidem illum sed reiciendis error, eligendi
            tenetur eos ullam, illo saepe deserunt!
          </p>

          <div bg-gray-50 p-4 flex flex-row-reverse gap-4>
            <button
              type="button"
              justify-center
              rounded-md
              bg-blue-6
              px-3
              py-2
              text-sm
              font-semibold
              text-white
              shadow-xs
              hover:bg-blue-5
              w-20
            >
              Save
            </button>
            <button
              type="button"
              justify-center
              rounded-md
              bg-white
              px-3
              py-2
              text-sm
              font-semibold
              text-gray-900
              ring-1
              shadow-xs
              ring-gray-300
              ring-inset
              hover:bg-gray-50
              w-20
            >
              Cancel
            </button>
          </div>
        </div>
      </slot>
    </div>
  </div>
</template>
