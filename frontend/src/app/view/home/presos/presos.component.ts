import { Component, OnInit } from '@angular/core';
import { PresoService, Preso } from '../../../service/preso/preso.service';
import { Nugget } from '../../../service/content/content.service';

@Component({
  selector: 'app-presos',
  templateUrl: './presos.component.html',
  styleUrls: ['./presos.component.scss']
})
export class PresosComponent implements OnInit {

  presos: Array<Nugget> = [];

  constructor(private presoService: PresoService) {} 

  getPresos(): void {
    this.presoService.getPresos().subscribe(
        presos => this.presos =  presos.nuggets.sort(compare).slice(-4, -1));
  } 

  ngOnInit(): void {
    this.getPresos();
  }

}

function compare(a:Nugget, b:Nugget) {
  if (a.timestamp < b.timestamp) {
    return -1;
  }
  if (a.timestamp > b.timestamp) {
    return 1;
  }
  // a must be equal to b
  return 0;
}
