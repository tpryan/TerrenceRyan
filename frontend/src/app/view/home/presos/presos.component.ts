import { Component, OnInit } from '@angular/core';
import { PresoService, Preso } from '../../../service/preso/preso.service';

@Component({
  selector: 'app-presos',
  templateUrl: './presos.component.html',
  styleUrls: ['./presos.component.scss']
})
export class PresosComponent implements OnInit {

  presos: Array<Preso> = [];

  constructor(private presoService: PresoService) {} 

  getPresos(): void {
    this.presoService.getPresos().subscribe(
        presos => this.presos =  presos.sort(compare).slice(-4, -1));
  } 

  ngOnInit(): void {
    this.getPresos();
  }

}

function compare(a:Preso, b:Preso) {
  if (a.published < b.published) {
    return -1;
  }
  if (a.published > b.published) {
    return 1;
  }
  // a must be equal to b
  return 0;
}
