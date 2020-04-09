import { Component, OnInit } from '@angular/core';
import { RepoService } from '../../../service/repo/repo.service';
import { Repo } from '../../../service/repo/repo';



@Component({
  selector: 'app-repos',
  templateUrl: './repos.component.html',
  styleUrls: ['./repos.component.scss']
})
export class ReposComponent implements OnInit {

  repos: Array<Repo> = [];

  constructor(private repoService: RepoService) {} 

  getRepos(): void {
    this.repoService.getRepos()
        .subscribe(repos => this.repos =  repos.slice(1, 6));
  }
    

  ngOnInit(): void {
    this.getRepos();
  }


}
