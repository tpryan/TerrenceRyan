import { Component, OnInit } from '@angular/core';
import { RepoService } from '../../../service/repo/repo.service';
import { Repo } from '../../../service/repo/repo';
import { Nugget } from '../../../service/content/content.service';



@Component({
  selector: 'app-repos',
  templateUrl: './repos.component.html',
  styleUrls: ['./repos.component.scss']
})
export class ReposComponent implements OnInit {

  repos: Array<Nugget> = [];

  constructor(private repoService: RepoService) {} 

  getRepos(): void {
    this.repoService.getRepos()
        .subscribe(repos => this.repos =  repos.nuggets.slice(1, 6));
  }
    

  ngOnInit(): void {
    this.getRepos();
  }


}
