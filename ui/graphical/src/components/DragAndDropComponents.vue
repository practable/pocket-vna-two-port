<template>
  <div class='container-fluid m-2'>
  <!-- Possible measurements -->
  <h4 class='txt-primary'>{{header}}</h4>
  <div class='d-flex flex-row justify-content-center'>
      <div v-for='image in display' :key='image.type' class='col-md-3'>
          <div class='dragbox' @dragstart="dragImage(image)" @dragover.prevent @dragenter.prevent @touchstart="dragImage(image); dropPort1()">
              <!-- <img class='dragbox-image' :id='image.type' :src='require(`/images/${image.img}.png`)' :alt='image.type' :hidden='port1 == image || port2 == image'> -->
              <img class='dragbox-image' :id='image.type' :src='"/images/" + image.img + ".png"' :alt='image.type' :hidden='port1 == image || port2 == image'>
          
            </div>
          <figcaption class='txt-primary fig'>{{ image.type }}</figcaption>
      </div>
      
  </div>


    <!-- Pocket VNA and ports -->
  <div class='row'>

     <!-- If a DUT (either 1,2,3,4) are connected then show the equivalent image and this connects both ports -->
    <div v-if="port1.type == 'dut1' || port2.type == 'dut1'" class='col-sm-12 pvna' @mousedown='removePort1(); removePort2()' @touchstart="removePort1(); removePort2()">
        <img class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-port-1-2-dut1.png' alt='pocket-vna-connected'>
      </div>

      <div v-else-if="port1.type == 'dut2' || port2.type == 'dut2'" class='col-sm-12 pvna' @mousedown='removePort1(); removePort2()' @touchstart="removePort1(); removePort2()">
        <img class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-2port-dut2.png' alt='pocket-vna-connected'>
      </div>

      <div v-else-if="port1.type == 'dut3' || port2.type == 'dut3'" class='col-sm-12 pvna' @mousedown='removePort1(); removePort2()' @touchstart="removePort1(); removePort2()">
        <img class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-2port-dut3.png' alt='pocket-vna-connected'>
      </div>

      <div v-else-if="port1.type == 'dut4' || port2.type == 'dut4'" class='col-sm-12 pvna' @mousedown='removePort1(); removePort2()' @touchstart="removePort1(); removePort2()">
        <img class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-2port-dut4.png' alt='pocket-vna-connected'>
      </div>

      <!-- If the thru is connected then similarly connect both and show image -->
      <div v-else-if="port1.type == 'thru' || port2.type == 'thru'" class='col-sm-12 pvna' @mousedown='removePort1(); removePort2()' @touchstart="removePort1(); removePort2()">
        <img class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-2port-thru.png' alt='pocket-vna-connected'>
      </div>

      <!-- Else check if cal standards are shown in either port separately -->
      <div v-else-if="port1.type == 'short' || port1.type == 'open' || port1.type == 'load'" class='col-sm-12 pvna' @mousedown='removePort1' @touchstart="removePort1">
        <img v-if="!syncPorts" class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-port-1-cal.png' alt='pocket-vna-connected'>
        <img v-else class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-port-1-2-cal.png' alt='pocket-vna-connected'>
      </div>

      <div v-else-if="port2.type == 'short' || port2.type == 'open' || port2.type == 'load'" class='col-sm-12 pvna' @mousedown='removePort2' @touchstart="removePort2">
        <img v-if="!syncPorts" class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-port-2-cal.png' alt='pocket-vna-connected'>
        <img v-else class='pvna-img' id='pvna-connected-dut-image' src='/images/pvna-connected-port-1-2-cal.png' alt='pocket-vna-connected'>
      </div>


        <!-- Else ports are empty -->
      <div v-else class='col-sm-12 pvna' @dragover.prevent @dragenter.prevent>
        <img class='pvna-img' id='pvna-image' src='/images/pvna-disconnected.png' alt='pocket-vna'>

            <div class='dropbox mb-2' id='port2' @drop='dropPort2' @dragstart='removePort2'>
                <!-- <img v-if='port2 != ""' class='dropbox-image' :id='port2.type' :src='"/images/" + port2.img + ".png"' :alt='port2.type'> -->
            </div>

            <div class='dropbox' id='port1' @drop='dropPort1' @dragstart='removePort1'>
                <!-- <img v-if='port1 != ""' class='dropbox-image' :id='port1.type' :src='"/images/" + port1.img + ".png"' :alt='port1.type'> -->
            </div>
      </div>

    

     
  </div>

</div>

</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
    name: 'DragAndDropComponents',
    props: ['header', 'display', 'syncPorts'],
    emit:['port1change', 'port2change'],
    data () {
        return {
            selected: '',
            port1: '',
            port2: ''
        }
    },
    components: {
        
    },
    computed:{
        ...mapGetters([
            'getSParams',
        ])
    },
    watch:{
        
    },
    created(){
        
    },
    mounted(){
        

    },
    methods:{
        ...mapActions([
            'setDraggable',
            
        ]),
        dragImage(img){
            this.selected = img;
        },
        dropPort1(){
            this.port1 = this.selected;
            this.$emit('port1change', this.port1);
        },
        removePort1(){
            this.port1 = '';
            this.$emit('port1change', '');
        },
        dropPort2(){
            this.port2 = this.selected;
            this.$emit('port2change', this.port2);
        },
        removePort2(){
            this.port2 = '';
            this.$emit('port2change', '');
        },
        dragTest(event){
            event.dataTransfer.effectAllowed = 'move';
            
            let element = event.target;
            if(element.classList.contains('drop-area')){
            console.log(element.id);
                event.dataTransfer.setData("text/html", element.id + "|" + element.childNodes[0].id);
                
            } else{
            while(element.parentNode){
                element = element.parentNode;
                console.log(element.id);
                if(element.classList.contains('drop-area')){
                    event.dataTransfer.setData("text/html", element.id + "|" + element.childNodes[0].id);
                    break;
                }
                }
            }
        },

        dropTest(event){
            event.preventDefault();
            event.stopPropagation();
            let dropData = event.dataTransfer.getData('text/html');
            let dropItems = dropData.split("|");
            let draggedZone = document.getElementById(dropItems[0]);
            let droppedElement = document.getElementById(event.target.id);
            let draggedID = dropItems[1];
            
            if(droppedElement != null && droppedElement.classList.contains('drop-area')){
                if(event.target.childNodes.length > 0){
                draggedZone.appendChild(event.target.childNodes[0]);
                }
                console.log(draggedID);
                droppedElement.appendChild(document.getElementById(draggedID));
            } 
            else if(droppedElement){
                let element = droppedElement;
                while(element.parentNode){
                element = element.parentNode;
                if(element.classList.contains('drop-area')){
                    console.log(element.childNodes[0]);
                    draggedZone.appendChild(element.childNodes[0]);
                    element.appendChild(document.getElementById(draggedID));
                    
                    break;
                }
                }
            }
            return false;
        }


    }
}
</script>

<style>

.dragbox{
    background-color: rgba(159, 159, 159, 0.75);
    border:2px dashed grey;
    width: 7.5vw;
    height:7.5vw;
}

.dragbox-image{
    max-height: 100%;
    max-width: 100%;
}

.dropbox{
    position: absolute;
    z-index: 10;
    background-color: rgba(159, 159, 159, 0.2);
    border:2px dashed grey;
    width: 5vw;
    height:5vw;
}

#port1{
    top: 57%;
    left: 65%;
}

#port2{
    top: 22%;
    left: 65%;
}

.dropbox-hidden{
    background-color: black;
    width: 5vw;
    height:5vw;
}

.dropbox-image{
    max-height: 100%;
    max-width: 100%;
}

.fig{
    width: 7.5vw;
}

.pvna{
    position:relative;
}

#pvna-image{
    max-width: 100%;
    max-height:100%;
    /*transform: rotate(-90deg);*/
}

#pvna-connected-cal-image{
    max-width: 100%;
    max-height:100%;
    /*transform: rotate(-90deg);*/
}

#pvna-connected-dut-image{
    max-width: 100%;
    max-height:100%;
    /*transform: rotate(-90deg);*/
}

</style>